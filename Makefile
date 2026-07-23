include .env
export

setup:
	go get -u github.com/swaggo/swag/cmd/swag
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g ./main.go -o ./docs
	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/files

build-docker:
	docker compose build --no-cache

run-local:
	docker compose up -d db postgres redis
	go run main.go

up:
	docker compose up

down:
	docker compose down

restart:
	docker compose restart

build:
	go build -v ./...

test:
	go test -v ./... -race -cover

seed:
	go run scripts/seeders/cmd/main.go up

seed-clear:
	go run scripts/seeders/cmd/main.go down

clean:
	docker compose down db postgres redis
	rm -rf .dbdata
