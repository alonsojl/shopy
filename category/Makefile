# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

##@ Development

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: linter
linter: ## Lint source code.
	@golangci-lint run -c .golangci.yml > linter.txt

.PHONY: swagger
swagger: ## Generate swagger documentation.
	@swag init --ot yaml,json -o ./assets -g ./lambda/main.go

.PHONY: download
download: ## Download Go dependencies.
	@go get github.com/aws/aws-lambda-go@v1.49.0
	@go get github.com/aws/aws-sdk-go-v2@v1.36.1
	@go get github.com/aws/aws-sdk-go-v2/config@v1.29.6
	@go get github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue@v1.18.4
	@go get github.com/aws/aws-sdk-go-v2/feature/s3/manager@v1.17.62
	@go get github.com/aws/aws-sdk-go-v2/service/dynamodb@v1.40.1
	@go get github.com/aws/aws-sdk-go-v2/service/s3@v1.77.0
	@go get github.com/go-ozzo/ozzo-validation/v4@v4.3.0
	@go get github.com/google/uuid@v1.6.0
	@go mod tidy

.PHONY: lambda
lambda: ## Build lambda function and compress it into a zip file.
	@rm -rf ./assets/lambda.zip ./assets/bootstrap
	@GOOS=linux GOARCH=arm64 go build -o ./assets/bootstrap ./lambda/*.go
	@zip -j ./assets/lambda.zip ./assets/bootstrap
