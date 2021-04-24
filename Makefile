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

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
