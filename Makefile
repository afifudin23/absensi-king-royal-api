ifneq (,$(wildcard .env))
include .env
export
endif

APP_NAME ?= $(APP_NAME)
DATABASE_URL ?= $(DATABASE_URL)

.PHONY: setup tidy run air build migrate-up migrate-down migrate-create migrate-version migrate-force

setup:
	go mod tidy
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
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
	@db_url="$(DATABASE_URL)"; \
	case "$$db_url" in \
		*@tcp\(*|*@unix\(*) ;; \
		*) db_url=$$(echo "$$db_url" | sed -E 's#@([^/]+:[0-9]+)/#@tcp(\1)/#');; \
	esac; \
	migrate -path migrations -database "$$db_url" up

migrate-down:
	@db_url="$(DATABASE_URL)"; \
	case "$$db_url" in \
		*@tcp\(*|*@unix\(*) ;; \
		*) db_url=$$(echo "$$db_url" | sed -E 's#@([^/]+:[0-9]+)/#@tcp(\1)/#');; \
	esac; \
	migrate -path migrations -database "$$db_url" down 1

migrate-create:
	@[ -n "$(name)" ] || (echo "name is required. Example: make migrate-create name=create_users_table" && exit 1)
	migrate create -ext sql -dir migrations -format "20060102150405" $(name)

migrate-version:
	@db_url="$(DATABASE_URL)"; \
	case "$$db_url" in \
		*@tcp\(*|*@unix\(*) ;; \
		*) db_url=$$(echo "$$db_url" | sed -E 's#@([^/]+:[0-9]+)/#@tcp(\1)/#');; \
	esac; \
	migrate -path migrations -database "$$db_url" version

migrate-force:
	@[ -n "$(version)" ] || (echo "version is required. Example: make migrate-force version=20260222055214" && exit 1)
	@db_url="$(DATABASE_URL)"; \
	case "$$db_url" in \
		*@tcp\(*|*@unix\(*) ;; \
		*) db_url=$$(echo "$$db_url" | sed -E 's#@([^/]+:[0-9]+)/#@tcp(\1)/#');; \
	esac; \
	migrate -path migrations -database "$$db_url" force $(version)
