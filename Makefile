BIN             = raptor
OUTPUT_DIR      = build
TEST_PROFILE    = testing
PROD_PROFILE    = production
DTEST_DIR       = deploy/terraform/$(TEST_PROFILE)
DPROD_DIR       = deploy/terraform/$(PROD_PROFILE)

export AWS_REGION = us-east-1

.PHONY: help
.DEFAULT_GOAL := help

build/linux: clean ## Build a linux binary ready to be zip'ed for AWS Lambda Deployment
	mkdir -p $(OUTPUT_DIR) && GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o $(OUTPUT_DIR)/$(BIN) .

build/release: build/linux ## Zip linux binary as AWS Deployment archive
	cd $(OUTPUT_DIR) && zip $(BIN).zip $(BIN)

deploy/testing: ## Deploy zip'ed archive to AWS testing account
	export AWS_PROFILE=$(TEST_PROFILE); cd $(DTEST_DIR) && terraform init && terraform apply

deploy/production: deploy/testing test/integration ## Deploy zip'ed archive to AWS production account
	export AWS_PROFILE=$(TEST_PROFILE); cd $(DPROD_DIR) && terraform init && terraform apply

clean: clean/linux ## Remove all build artifacts

clean/linux: ## Remove linux build artifacts
	$(RM) $(OUTPUT_DIR)/$(BIN).zip
	$(RM) $(OUTPUT_DIR)/$(BIN)

test/integration: ## Integration Testing
	AWS_PROFILE=$(TEST_PROFILE) go test -tags integration -timeout 30s ./job -run ^TestS3Events$$ -v

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'	