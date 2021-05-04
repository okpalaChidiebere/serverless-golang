.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/src/lambda/http/getGroups src/lambda/http/getGroups/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/src/lambda/http/createGroup src/lambda/http/createGroup/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getImages src/lambda/http/getImages/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/getImage src/lambda/http/getImage/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/createImage src/lambda/http/createImage/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/sendNotifications src/lambda/s3/sendNotifications/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/connect src/lambda/websocket/connect/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/disconnect src/lambda/websocket/disconnect/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/elasticSearchSync src/lambda/dynamoDb/elasticSearchSync/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/resizeImage src/lambda/s3/resizeImage/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/auth0Authorizer src/lambda/auth/auth0Authorizer/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/src/models/Group src/models/Group.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/src/requests/CreateGroupRequest src/requests/CreateGroupRequest.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/src/businessLogic/groups/groups src/businessLogic/groups/groups.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/src/dataLayer/groupsAccess/groupsAccess src/dataLayer/groupsAccess/groupsAccess.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
