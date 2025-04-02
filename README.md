# Serverless GoLang

This project was about learning how to learn to write a serverless application in GoLang. Every code writtn here was already done in NodeJs [here](https://github.com/okpalaChidiebere/cloud-developer/tree/master/course-04/exercises/lesson-6/solution). I had to learn to write them in Golang!

To bootstrap an AWS Go Serverless template follow there are three types of template

- aws-go for basic services

- aws-go-mod uses go modules

- aws-go-dep used go dep. You will have an aws-sdk vendor folder

I prefer the aws-go-mod template. To get started

- Make sure you have [Serverless](https://www.serverless.com/framework/docs/getting-started) CLI installed or upgraded. To Install or Upgrade CLI run `npm install -g serverless`
- create folder for your app and name it whatever
- `cd` into that folder
- inside that folder run `sls create --template aws-go-mod`
- **Note** With the deprecation of `go1.x` scheduled for March, 2024, it leaves us the provided family of runtimes as our only option. In our case, we will be using `provided.al2`, which is the latest available runtime, based on Amazon Linux 2023. When using provided runtimes, there is a requirement to name the executable for your function bootstrap. See [here](https://blog.matthiasbruns.com/running-multiple-golang-aws-lambda-functions-on-arm64-with-serverlesscom),and [here](https://pgrzesik.com/posts/golang-serverless-go-plugin/). Another article ARM64 lambda golang [here](https://blog.matthiasbruns.com/running-multiple-golang-aws-lambda-functions-on-arm64-with-serverlesscom)

```makefile
build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
```

becomes

```makefile
build:
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/hello/bootstrap hello/main.go

zip:
	zip -j build/hello.zip build/hello/bootstrap
```

```yml
provider:
  name: aws
  runtime: go1.x

functions:
  hello:
    handler: bin/hello
    events:
      - httpApi:
          path: /hello
          method: get
```

becomes

```yml
provider:
  name: aws
  runtime: provided.al2
  architecture: arm64 # Lambda binary must be is compiled for arm64 in your Makefile and not amd64

package:
  individually: true # <- package each function individually, to prevent file name conflicts

functions:
  hello:
    handler: bootstrap # <- the handler name must be bootstrap and in the root of the zip
    package:
      artifact: build/hello.zip
    events:
      - httpApi:
          path: /hello
          method: get
```

Then after, to have your first deploy to aws

- Make sure you have an aws user logged in in your CLI. Make your you have set up [aws-iam-authenticator](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html). To login a user run `aws configure` and enter user credentials. To confirm the user that you have locally login into your awscli run `aws sts get-caller-identity` in terminal
- Use the make file to build your project. In the terminal just run `make`. NOTE you probably will have to install some dependencies(like aws-sdk-go, aws-lambda-go, etc) for the make file to successfully generate executable for your project with `go mod init <modulename>` and `go mod tidy`
- once your executables are generated, run `sls deploy --verbose` to deploy your project to aws
  All the steps i listed is in this article [https://schadokar.dev/posts/create-a-serverless-application-in-golang-with-aws/](https://schadokar.dev/posts/create-a-serverless-application-in-golang-with-aws/)

# More related articles on bootstrapping a sls go template

[https://tpaschalis.github.io/golang-aws-lambda-getting-started/](https://tpaschalis.github.io/golang-aws-lambda-getting-started/)
[https://www.softkraft.co/aws-lambda-in-golang/](https://www.softkraft.co/aws-lambda-in-golang/)

# Articles on writing GoLang code for lambda

- [https://yos.io/2018/02/08/getting-started-with-serverless-go/](https://yos.io/2018/02/08/getting-started-with-serverless-go/)
- [https://dev.to/jeastham1993/aws-dynamodb-in-golang-1m2j](https://dev.to/jeastham1993/aws-dynamodb-in-golang-1m2j)
- [https://github.com/packtpublishing/hands-on-serverless-applications-with-go](https://github.com/packtpublishing/hands-on-serverless-applications-with-go)

# Interesting Read about optimizing golang runtime and bills

- [https://forum.serverless.com/t/optimizing-lambdas-reducing-your-bills/4101/2]https://forum.serverless.com/t/optimizing-lambdas-reducing-your-bills/4101/2
- [https://runbook.cloud/blog/posts/how-we-massively-reduced-our-aws-lambda-bill-with-go/](https://runbook.cloud/blog/posts/how-we-massively-reduced-our-aws-lambda-bill-with-go/)
- [https://www.simplybusiness.co.uk/about-us/tech/2021/02/go-routines-aws-lambda/](https://www.simplybusiness.co.uk/about-us/tech/2021/02/go-routines-aws-lambda/)

# Go Lambda Middlewares

Right now there is no really good middleware to use right now but there are some third party libraries will soon be the standard ones just like middy for NodeJS aws lambda. Check the link below

- [https://github.com/mefellows/vesper](https://github.com/mefellows/vesper)
- [https://jpcedeno.com/post/gointercept-introduction/](https://jpcedeno.com/post/gointercept-introduction/)

In the AWS X-ray i implemented for go, i keep getting a context error, i can resolve it later here

- This [link](https://medium.com/nordcloud-engineering/tracing-serverless-application-with-aws-x-ray-2b5e1a9e9447) is a more complete tutorial on how to fully properly used the X-ray in go
- [https://www.gitmemory.com/issue/aws/aws-xray-sdk-go/50/548321975](https://www.gitmemory.com/issue/aws/aws-xray-sdk-go/50/548321975)
- [https://github.com/aws/aws-xray-sdk-go/issues/50](https://github.com/aws/aws-xray-sdk-go/issues/50)

Full documentation of aws-sdk-go [here](https://docs.aws.amazon.com/sdk-for-go/api/aws/)
More links

- [https://github.com/aws/aws-lambda-go/tree/master/events](https://github.com/aws/aws-lambda-go/tree/master/events)

File Upload articles to AWS S3 in go
[https://medium.com/spankie/upload-images-to-aws-s3-bucket-in-a-golang-web-application-2612bea70dd8](https://medium.com/spankie/upload-images-to-aws-s3-bucket-in-a-golang-web-application-2612bea70dd8)
[https://stackoverflow.com/questions/49266516/reading-files-from-aws-s3-in-golang](https://stackoverflow.com/questions/49266516/reading-files-from-aws-s3-in-golang)
[https://questhenkart.medium.com/s3-image-uploads-via-aws-sdk-with-golang-63422857c548](https://questhenkart.medium.com/s3-image-uploads-via-aws-sdk-with-golang-63422857c548)

PhotoShop library in go
[https://github.com/anthonynsimon/bild](https://github.com/anthonynsimon/bild)

# API Gateway Stages

- We can have different stages like `dev`, `staging` and/or `prod`
- Every on of thses stages will have a different url. Eg `api-gateway.com/dev`, `api-gateway.com/staging`, `api-gateway.com/prod`
- These differnt stages will use different versions of our Lambda functions. The dev will have the most recent version and prod will have a older version. Eg prod will have version 4, staging version 7 and prod version 10
- I did not implement this, but this is the ideal way to stage your gateway in aws. I will look for a way to do this with serverless framework

# HTTP API (API Gateway v2)

- v1, also called REST API which is what we used for this project
- v2, also called [HTTP API](https://www.serverless.com/framework/docs/providers/aws/events/http-api/), which is faster and cheaper than v1. Read the full comparison [in the AWS documentation](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-vs-rest.html).

- Why you may need to [https://www.serverless.com/blog/api-gateway-v2-http-apis](https://www.serverless.com/blog/api-gateway-v2-http-apis)

- How to do so in serverless [https://www.serverlessguru.com/blog/migrating-from-api-gateway-v1-rest-to-api-gateway-v2-http](https://www.serverlessguru.com/blog/migrating-from-api-gateway-v1-rest-to-api-gateway-v2-http)

- [https://www.serverless.com/aws-http-apis](https://www.serverless.com/aws-http-apis)

- More about CORS errors [here](https://repost.aws/knowledge-center/api-gateway-cors-errors) and [here](https://www.stackhawk.com/blog/golang-cors-guide-what-it-is-and-how-to-enable-it/)

A refresher for common basic goLang programming techniques for you [here](https://www.bogotobogo.com/GoLang/GoLang_Modules_1_Creating_a_new_module.php)

# Auth0 Authentication and jwt token helpful links

- [Implementing refresh token flow in an expo react native app with expo-auth-session and Auth0](https://medium.com/@danbowden/implementing-refresh-token-flow-in-an-expo-react-native-app-with-expo-auth-session-and-auth0-82eb6d0dea35) and [here](https://gist.github.com/jdthorpe/aaa0d31a598f299a57e5c76535bf0690)
- [Auth0 discovery url](https://community.auth0.com/t/openid-discovery-url/19536)
- Get test accessToken. [see](https://auth0.com/docs/secure/tokens/access-tokens/get-management-api-access-tokens-for-testing)
- AccessToken claims. [See](https://auth0.com/docs/secure/tokens/access-tokens)
- [Auth0 Public Key for token verifications (Asymmetric)](https://community.auth0.com/t/how-to-get-public-key-pem-from-jwks-json/60355/6)
- [Client Secret for token verifications (Symmetric)](https://auth0.com/docs/get-started/tenant-settings/signing-keys/view-signing-certificates#if-using-the-hs256-signing-algorithm)
- [Public key for Asymmetric token verification](https://auth0.com/docs/secure/tokens/json-web-tokens/validate-json-web-tokens#verify-rs256-signed-tokens)
- Parse jwt token without verification to get claims. see [here](https://stackoverflow.com/questions/45405626/how-to-decode-a-jwt-token-in-go) and [here](https://stackoverflow.com/questions/55698770/decode-jwt-without-validation-and-find-scope)
- [Lambda Authorizer caching](https://stackoverflow.com/questions/50331588/aws-api-gateway-custom-authorizer-strange-showing-error/56119016#56119016)

# Auth0 and API authorization

Authorization is basically making sure that a user that invoke an endpoint has the permission to assess that resource. With Auth0 you can achieve this using Actions. Follow the steps below

- Create an API Audience. Go to to Dashboard > Applications > APIs. You must provide this audience as params following auth0 sign in otherwise verifying the access token following the Auth0 verification process will fail. As you create the Audience, create all the permissions for the API client
- Create the role(s) at Go to to Dashboard > User management > Roles. While creating the role, you will assign the permissions you want the role to have from the API audience you created.
- You can manually assign role to user(s). You can see all users at Dashboard > User management > Users
- For cases where you want to assign default role to users when an Auth0 event(flow) happens like Login/Post signUp which is a common use case, you will need to create an action. Go to Dashboard > Actions > Library and click Create Action / Build from scratch. When you are done creating the action, then you can hook the action up into a flow (Dashboard > Actions > Flows)
- **Note** It is important that the client you will be using to assign role to users has the correct permission scope assign to then in the auth0 dashboard. See. The minimum permission needed for this is `and`
- [https://auth0.com/blog/assign-default-role-on-sign-up-with-actions/](https://auth0.com/blog/assign-default-role-on-sign-up-with-actions/)
- [https://auth0.com/docs/customize/actions/flows-and-triggers](https://auth0.com/docs/customize/actions/flows-and-triggers)
- [https://www.youtube.com/watch?v=CZxfMD8lXg8](https://www.youtube.com/watch?v=CZxfMD8lXg8)
- [Scopes and Claims](https://auth0.com/docs/get-started/apis/scopes/sample-use-cases-scopes-and-claims)
- Enable Role-Based Access Control for APIs. See [Here](https://auth0.com/docs/get-started/apis/enable-role-based-access-control-for-apis), [here](https://community.auth0.com/t/how-to-get-permissions-for-user/26993/10) and [here](https://community.auth0.com/t/get-auth0-management-api-invalid-token-error/83649/8)
- [https://community.auth0.com/t/how-to-get-permissions-for-user/26993/13](https://community.auth0.com/t/how-to-get-permissions-for-user/26993/13)

## Auth0 APis front-end useful links

- [Prompt parameter](https://community.auth0.com/t/authorizes-prompt-parameter-not-documented/40340/2)

## DynamoDB Links

- For TTL, see [here](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/time-to-live-ttl-how-to.html) and [here](https://dynobase.dev/dynamodb-ttl/)
- [https://www.bmc.com/blogs/dynamodb-advanced-queries/](https://www.bmc.com/blogs/dynamodb-advanced-queries/)
- [https://www.dynamodbguide.com/working-with-multiple-items/](https://www.dynamodbguide.com/working-with-multiple-items/)
