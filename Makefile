GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
BINARY_NAME=HARnonymize
BINARY_DIR=jfallis/harnonymize
VERSION?=1.0.0
EXPORT_RESULT?=false

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build vendor

all: help

## Build:
vendor: ## Install the dependencies
	$(GOCMD) mod vendor

build: vendor ## Build the project and put the output binary
	mkdir -p out/bin
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_DIR)/$(BINARY_NAME)-windows.exe
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_DIR)/$(BINARY_NAME)-amd64-macos
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_DIR)/$(BINARY_NAME)-arm64-macos
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_DIR)/$(BINARY_NAME)-linux
	cp block.txt $(BINARY_DIR)/block.txt
	zip -r $(BINARY_NAME).zip $(BINARY_DIR)

clean: ## Remove build related files
	rm -fr ./bin
	rm -fr ./out
	rm -f ./junit-report.xml checkstyle-report.xml ./coverage.xml ./profile.cov yamllint-checkstyle.xml

## Test:
lint: vendor ## golang linting
	golangci-lint run

test: vendor ## Run the tests
ifeq ($(EXPORT_RESULT), true)
	go get -u github.com/jstemmer/go-junit-report
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	$(GOTEST) -v -race ./... $(OUTPUT_OPTIONS)

bench: vendor ## Run the benchmarks
	$(GOTEST) -v -bench=. ./...

coverage: vendor ## Run the tests and export the coverage
	$(GOTEST) -cover -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
