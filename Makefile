ifneq ("$(wildcard .env)","")
	include .env
	export $(shell sed 's/=.*//' .env)
endif

.PHONY: start
start:
	./bin/main

.PHONY: build
build:
	go build -o=./bin/main ./cmd

.PHONY: test
test:
	go test -v ./...
	
.PHONY: codegen
codegen:
	go run github.com/99designs/gqlgen

# https://github.com/golang-migrate/migrate

.PHONY: models
models:
	pg_dump --schema-only go_api_seed > schema.sql
	sqlc generate

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrate-up
migrate-up:
	migrate -path=./migrations -database="$(POSTGRES_URL)" up

.PHONY: migrate-down
migrate-down:
	migrate -path=./migrations -database="$(POSTGRES_URL)" down 1

.PHONY: test-env
test-env:
	@echo $(POSTGRES_URL)