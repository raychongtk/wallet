.PHONY: pre-commit start stop

pre-commit:
	go mod tidy
	go mod vendor
	wire
	go vet
	go fmt ./...

start:
	docker-compose up -d
	go run github.com/raychongtk/wallet

stop:
	docker-compose down