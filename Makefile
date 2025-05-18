.PHONY: pre-commit start stop test

pre-commit:
	go mod tidy
	go mod vendor
	wire
	go vet
	go fmt ./...

start:
	docker-compose up -d
	sleep 5
	go run github.com/raychongtk/wallet

stop:
	docker-compose down

test:
	go test -json ./...