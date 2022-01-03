.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	sudo docker-compose up --remove-orphans --build app

swag:
	swag init -g cmd/app/main.go

deploy:
	sudo docker-compose -f docker-compose-deploy.yaml up --remove-orphans --build