# Serverless GoLang

This project was about learning how to learn to write a serverless application in GoLang. Every code writtn here was already done in NodeJs [here](https://github.com/okpalaChidiebere/cloud-developer/tree/master/course-04/exercises/lesson-6/solution). I had to learn to write them in Golang!

To bootstarap an AWS Go Serverless template follow  there are three types of template 
* aws-go for basic services
* aws-go-mod uses go modules
* aws-go-dep used go dep. You will have an aws-sdk vendor folder

I prefer the aws-go-mod template. To get started
* create folder for your app and name it whatever
* `cd` into that folder
* inside that folder run `sls create --template aws-go-mod`

Then after to have your first deploy
* use the make file to build your project. In the terminal just run `make`. NOTE you probably will have to install some dependencies(like aws-sdk-go, aws-lambda-go, etc) for the make file to successfully generate executable for your project 
* once your executables are generated, run `sls deploy -v` to deploy your project to  aws
All the steps i listed is in this article [https://schadokar.dev/posts/create-a-serverless-application-in-golang-with-aws/](https://schadokar.dev/posts/create-a-serverless-application-in-golang-with-aws/)

# More related articles on bootstraping a sls go template
[https://tpaschalis.github.io/golang-aws-lambda-getting-started/](https://tpaschalis.github.io/golang-aws-lambda-getting-started/)
[https://www.softkraft.co/aws-lambda-in-golang/](https://www.softkraft.co/aws-lambda-in-golang/)

# Articles on writing GoLang code for lambda
* [https://yos.io/2018/02/08/getting-started-with-serverless-go/](https://yos.io/2018/02/08/getting-started-with-serverless-go/)
* [https://dev.to/jeastham1993/aws-dynamodb-in-golang-1m2j](https://dev.to/jeastham1993/aws-dynamodb-in-golang-1m2j)
* [https://github.com/packtpublishing/hands-on-serverless-applications-with-go](https://github.com/packtpublishing/hands-on-serverless-applications-with-go
)

# Interesting Read about optimizing golang runtime and bills
* [https://forum.serverless.com/t/optimizing-lambdas-reducing-your-bills/4101/2]https://forum.serverless.com/t/optimizing-lambdas-reducing-your-bills/4101/2
* [https://runbook.cloud/blog/posts/how-we-massively-reduced-our-aws-lambda-bill-with-go/](https://runbook.cloud/blog/posts/how-we-massively-reduced-our-aws-lambda-bill-with-go/)
* [https://www.simplybusiness.co.uk/about-us/tech/2021/02/go-routines-aws-lambda/](https://www.simplybusiness.co.uk/about-us/tech/2021/02/go-routines-aws-lambda/)

# Go Lambda Middlewares
Right now there is no really good middleware to use right now but there are some third party libraries will soon be the standard ones just like middy for NodeJS aws lambda. Check the link below
* [https://github.com/mefellows/vesper](https://github.com/mefellows/vesper)
* [https://jpcedeno.com/post/gointercept-introduction/](https://jpcedeno.com/post/gointercept-introduction/)

In the AWS X-ray i implemented for go, i keep getting a context error, i can resolve it later here
* This [link](https://medium.com/nordcloud-engineering/tracing-serverless-application-with-aws-x-ray-2b5e1a9e9447) is a more complete tutorial on how to fully properly used the X-ray in go 
* [https://www.gitmemory.com/issue/aws/aws-xray-sdk-go/50/548321975](https://www.gitmemory.com/issue/aws/aws-xray-sdk-go/50/548321975)
* [https://github.com/aws/aws-xray-sdk-go/issues/50](https://github.com/aws/aws-xray-sdk-go/issues/50)


Full documentation of aws-sdk-go [here](https://docs.aws.amazon.com/sdk-for-go/api/aws/)
More links
* [https://github.com/aws/aws-lambda-go/tree/master/events](https://github.com/aws/aws-lambda-go/tree/master/events)


File Upload articles to AWS S3 in go
[https://medium.com/spankie/upload-images-to-aws-s3-bucket-in-a-golang-web-application-2612bea70dd8](https://medium.com/spankie/upload-images-to-aws-s3-bucket-in-a-golang-web-application-2612bea70dd8)
[https://stackoverflow.com/questions/49266516/reading-files-from-aws-s3-in-golang](https://stackoverflow.com/questions/49266516/reading-files-from-aws-s3-in-golang)
[https://questhenkart.medium.com/s3-image-uploads-via-aws-sdk-with-golang-63422857c548](https://questhenkart.medium.com/s3-image-uploads-via-aws-sdk-with-golang-63422857c548)


PhotoShop library in go
[https://github.com/anthonynsimon/bild](https://github.com/anthonynsimon/bild)

# API Gateway Stages
* We can have different stages like `dev`, `staging` and/or `prod`
* Every on of thses stages will have a different url. Eg `api-gateway.com/dev`, `api-gateway.com/staging`, `api-gateway.com/prod`
* These differnt stages will use different versions of our Lambda functions. The dev will have the most recent version and prod will have a older version. Eg prod will have version 4, staging version 7 and prod version 10
* I did not implement this, but this is the ideal way to stage your gateway in aws. I will look for a way to do this with serverless framework


Keep an EYE out for APIGatway [Version2](https://www.serverless.com/framework/docs/providers/aws/events/http-api/)
* Why you may need to [https://www.serverless.com/blog/api-gateway-v2-http-apis](https://www.serverless.com/blog/api-gateway-v2-http-apis)
* How to do so in serverless [https://www.serverlessguru.com/blog/migrating-from-api-gateway-v1-rest-to-api-gateway-v2-http](https://www.serverlessguru.com/blog/migrating-from-api-gateway-v1-rest-to-api-gateway-v2-http)
* [https://www.serverless.com/aws-http-apis](https://www.serverless.com/aws-http-apis)

A refresher for common basic goLang programming techniques for you [here](https://www.bogotobogo.com/GoLang/GoLang_Modules_1_Creating_a_new_module.php)


