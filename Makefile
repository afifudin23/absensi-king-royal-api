SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command

ifneq (,$(wildcard .env))
include .env
export
endif

APP_NAME ?= $(APP_NAME)
DATABASE_URL ?= $(DATABASE_URL)

.PHONY: setup tidy run air build swagger migrate-up migrate-down migrate-create migrate-version migrate-force

setup:
	go mod tidy
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest

swagger:
	swag init -g cmd/api/main.go --output docs --parseDependency --parseInternal

tidy:
	go mod tidy

run:
	go run ./cmd/api

air:
	air

build:
	go build -o bin/$(APP_NAME) ./cmd/api

migrate-up:
	@$$url = '$(DATABASE_URL)'; if ($$url -notmatch '@tcp\(|@unix\(') { $$url = $$url -replace '@([^/]+:\d+)/', '@tcp($$1)/' }; migrate -path migrations -database $$url up

migrate-down:
	@$$url = '$(DATABASE_URL)'; if ($$url -notmatch '@tcp\(|@unix\(') { $$url = $$url -replace '@([^/]+:\d+)/', '@tcp($$1)/' }; migrate -path migrations -database $$url down 1

migrate-create:
	@if (-not '$(name)') { throw 'name is required. Example: make migrate-create name=create_users_table' }; migrate create -ext sql -dir migrations -format '20060102150405' $(name)

migrate-version:
	@$$url = '$(DATABASE_URL)'; if ($$url -notmatch '@tcp\(|@unix\(') { $$url = $$url -replace '@([^/]+:\d+)/', '@tcp($$1)/' }; migrate -path migrations -database $$url version

migrate-force:
	@if (-not '$(version)') { throw 'version is required. Example: make migrate-force version=20260222055214' }; $$url = '$(DATABASE_URL)'; if ($$url -notmatch '@tcp\(|@unix\(') { $$url = $$url -replace '@([^/]+:\d+)/', '@tcp($$1)/' }; migrate -path migrations -database $$url force $(version)