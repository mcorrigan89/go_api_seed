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
	migrate -path=./migrations -database="$(POSTGRES_URL)?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	migrate -path=./migrations -database="$(POSTGRES_URL)?sslmode=disable" down 1

.PHONY: test-env
test-env:
	@echo $(POSTGRES_URL)

.PHONY: install-new-relic
install-new-relic:
	curl -Ls https://download.newrelic.com/install/newrelic-cli/scripts/install.sh | bash && sudo NEW_RELIC_API_KEY=$(NEW_RELIC_API_KEY) NEW_RELIC_ACCOUNT_ID=$(NEW_RELIC_ACCOUNT_ID) /usr/local/bin/newrelic install -n logs-integration
