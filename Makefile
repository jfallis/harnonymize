BINARY_NAME=harnonymize

test:
	go test --race -v ./...

test-lint:
	golangci-lint run ./...

build:
	GOOS=windows GOARCH=amd64 go build -o HARnonymize/$(BINARY_NAME)-windows.exe
	GOOS=darwin GOARCH=amd64 go build -o HARnonymize/$(BINARY_NAME)-amd64-macos
	GOOS=darwin GOARCH=arm64 go build -o HARnonymize/$(BINARY_NAME)-arm64-macos
	GOOS=linux GOARCH=amd64 go build -o HARnonymize/$(BINARY_NAME)-linux
	cp block.txt HARnonymize/block.txt
	zip -r HARnonymize.zip HARnonymize