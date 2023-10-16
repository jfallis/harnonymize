BINARY_NAME=harnonymize

test:
	go test --race -v ./...

test-lint:
	golangci-lint run ./...

build:
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows.exe
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-amd64-macos
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-arm64-macos
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux
