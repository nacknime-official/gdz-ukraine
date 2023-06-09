include .env
export

BINDIR    := $(CURDIR)/bin
PATH      := $(BINDIR):$(PATH)
GOFLAGS	  := -trimpath

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: ### build
	go build $(GOFLAGS) -o $(BINDIR)/ ./cmd/app

run: ### run
	go run ./cmd/app

test: ### run test
	go test -v -cover -race ./...
.PHONY: test

gen: ### go generate
	go generate ./...
.PHONY: gen

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations -seq -digits 4 'migrate_name'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(POSTGRES_URL)' up
.PHONY: migrate-up

migrate-clean: ### migration clean database
	migrate -path migrations -database '$(POSTGRES_URL)' drop -f
.PHONY: migrate-up

migrate-version: ### migration version
	migrate -path migrations -database '$(POSTGRES_URL)' version
.PHONY: migrate-up

bin-deps: ### install binary dependencies (gen tools, migrator etc)
	GOBIN=$(BINDIR) go install -tags 'pgx' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2
	GOBIN=$(BINDIR) go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.17.2

