##@ General

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

lambdas := category product user

## Display this help.
.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

## Start container services.
.PHONY: up
up:
	@docker compose up -d
	
## Stop container services.
.PHONY: down
down:
	@docker compose down

## Deploy the stack into your AWS account.
.PHONY: deploy
deploy:
	@cdk deploy

## Destroy the stack from your AWS account.
.PHONY: destroy
destroy:
	@cdk destroy

.PHONY: lambdas
lambdas: ## Build lambda functions.
	@$(foreach lambda, $(lambdas), (cd $(lambda) && $(MAKE) lambda);)

## Run & test AWS serverless functions locally as a HTTP API.
.PHONY: local
local: template
	@sam local start-api --template cdk.out/ShopyStack.template.json

## Generate a CloudFormation template in YAML format.
.PHONY: template
template:
	@cdk synth > template.yaml
