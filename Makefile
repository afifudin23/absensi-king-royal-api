ifneq (,$(wildcard .env))
include .env
export
endif

APP_NAME ?= $(APP_NAME)
DATABASE_URL ?= $(DATABASE_URL)

.PHONY: setup tidy run air build migrate-up migrate-down migrate-create

setup:
	go mod tidy
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/air-verse/air@latest

tidy:
	go mod tidy

run:
	go run ./cmd/api

air:
	air

build:
	go build -o bin/$(APP_NAME) ./cmd/api

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	@[ -n "$(name)" ] || (echo "name is required. Example: make migrate-create name=create_users_table" && exit 1)
	migrate create -ext sql -dir migrations -format "20060102150405" $(name)
