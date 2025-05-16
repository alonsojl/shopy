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

##@ Development

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: up
up: ## Start container services.
	@docker compose up -d
	
.PHONY: down
down: ## Stop container services.
	@docker compose down

.PHONY: deploy
deploy: ## Deploy the stack into your AWS account.
	@cdk deploy

.PHONY: destroy
destroy: ## Destroy the stack from your AWS account.
	@cdk destroy

.PHONY: lambdas
lambdas: ## Build lambda functions.
	@$(foreach lambda, $(lambdas), (cd $(lambda) && $(MAKE) lambda);)

.PHONY: local
local: template ## Run & test AWS serverless functions locally as a HTTP API.
	@sam local start-api --template cdk.out/ShopyStack.template.json

.PHONY: template
template: ## Generate a CloudFormation template in YAML format.
	@cdk synth > template.yaml
