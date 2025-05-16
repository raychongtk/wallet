.PHONY: pre-commit

pre-commit:
	go mod tidy
	go mod vendor
	wire
	go vet
	go fmt ./...