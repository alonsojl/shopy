package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

// @title 			Product Lambda Function
// @version 		1.0
// @description 	Testing Swagger APIs.
// @termsOfService 	http://swagger.io/terms/

// @schemes 		http https
// @host            127.0.0.1:3000
// @basePath 		/dev/v1

// @contact.name 	API Support
// @contact.url		http://www.swagger.io/support
// @contact.email 	support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// @securityDefinitions.apiKey JWT
// @in	 header
// @name Authorization

func main() {
	lambda.Start(handler.Router())
}
