.PHONY: pre-commit run

pre-commit:
	go mod tidy
	go mod vendor
	wire
	go vet
	go fmt ./...

run:
	docker-compose up -d
	go run github.com/raychongtk/wallet

stop:
	docker-compose down