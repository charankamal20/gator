run:
	@go run .

test:
	@go test -v ./...

lint:
	@go run golang.org/x/lint/golint ./...

build:
	@go build -o bin/ .
