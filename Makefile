.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getGroups src/lambda/http/getGroups/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/createGroup src/lambda/http/createGroup/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getImages src/lambda/http/getImages/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getImage src/lambda/http/getImage/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/createImage src/lambda/http/createImage/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/sendNotifications src/lambda/s3/sendNotifications/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/connect src/lambda/websocket/connect/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/disconnect src/lambda/websocket/disconnect/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/elasticSearchSync src/lambda/dynamoDb/elasticSearchSync/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/resizeImage src/lambda/s3/resizeImage/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
