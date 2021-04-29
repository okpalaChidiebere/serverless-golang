package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nfnt/resize"
)

type s3Event events.S3Event
type SNSEvent events.SNSEvent

var (
	thumbnailBucketName = os.Getenv("THUMBNAILS_S3_BUCKET")
	imagesBucketName    = os.Getenv("IMAGES_S3_BUCKET")
	s3Client            *s3.S3
)

func init() {
	svc := session.Must(session.NewSession())
	s3Client = s3.New(svc)
}

func main() {
	lambda.Start(resizeImageHandler)
}

func resizeImageHandler(e SNSEvent) {
	for _, record := range e.Records {
		var s3Event s3Event
		c := make(chan string)
		snsRecord := record.SNS

		fmt.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)

		//we parse the message to get the actual S3Event
		if err := json.Unmarshal([]byte(snsRecord.Message), &s3Event); err != nil {
			log.Println("Failed to unmarshal")
		}

		for _, s3EventRecord := range s3Event.Records {
			go processImage(s3EventRecord, c)
		}

		for i := 0; i < len(s3Event.Records); i++ {
			fmt.Printf("S3 record done. Key = %s", <-c) //each time this print line is called, it block the main thread from execting, so we can still get all our results back
		}
	}
}

func processImage(e events.S3EventRecord, c chan string) {
	key := e.S3.Object.Key
	fmt.Printf("Processing S3 item with key: %s", key)

	req, resp := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(imagesBucketName),
		Key:    aws.String(key),
	})

	if err := req.Send(); err != nil {
		fmt.Print(err)
		c <- key
		return
	}

	s3ObjectBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
		c <- key
		return
	}

	fmt.Print("Resizing image")
	//reset format of data []byte to image.Image
	img, _, err := image.Decode(bytes.NewReader(s3ObjectBytes))
	if err != nil {
		fmt.Print(err)
		c <- key
		return
	}
	newImage := resize.Resize(150, 0, img, resize.Lanczos2)
	buf := new(bytes.Buffer)
	//convert our image.Image into a buffer
	if err := jpeg.Encode(buf, newImage, nil); err != nil {
		fmt.Print(err)
		c <- key
		return
	}

	fmt.Printf("Writing image back to S3 bucket: %s", thumbnailBucketName)
	//Uploading the resized image to another S3 bucket
	res, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(thumbnailBucketName),
		Key:    aws.String(key + ".jpeg"),
		Body:   bytes.NewReader(buf.Bytes()), //provide the buffer we want to write to this S3 bucket
	})
	if err != nil {
		fmt.Printf("failed to write image back to bucket. Error: %s", err)
		c <- key
		return
	}

	fmt.Print(res.String() + "\n")
	c <- key

}
