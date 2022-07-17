include .env
export

up:
	docker-compose up --build -d

test: db
	go test ./...

db:
	docker-compose up --build -d db 

swagger:
	swag init -g cmd/main.go

run:
	go run cmd/main.go