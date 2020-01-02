SLACK_TOKEN :=xoxo-

.PHONY: build
## build: build the application
build: clean
	@echo "Building..."
	@go build -o ${APP} cmd/brownie/main.go

.PHONY: run
## run: runs go run main.go
run:
	go run -race cmd/brownie/main.go

.PHONY: clean
## clean: cleans the binary
clean:
	@echo "Cleaning"
	@go clean

.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...

.PHONY: build-tokenizer
## build-tokenizer: build the tokenizer application
build-tokenizer:
	${MAKE} -c tokenizer build

.PHONY: setup
## setup: setup go modules
setup:
	@go mod init \
		&& go mod tidy \
		&& go mod vendor
