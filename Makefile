SOURCES := $(shell find . -name '*.go' -type f -not -path './vendor/*'  -not -path '*/mocks/*')

PACKAGE := github.com/uudashr/jobexec
OBJ_DIR := $(GOPATH)/pkg/darwin_amd64/$(PACKAGE)

# Database
DB_USER ?= jobexec
DB_PASSWORD ?= jobexecsecret
DB_PORT ?= 3306
DB_ADDRESS ?= 127.0.0.1:${DB_PORT}
DB_NAME ?= jobexec_test

# Dependencies Management
.PHONY: vendor-prepare
vendor-prepare:
	@echo "Installing dep"
	@go get -u github.com/golang/dep/cmd/dep

Gopkg.lock: Gopkg.toml
	@dep ensure -update

.PHONY: vendor-update
vendor-update:
	@dep ensure -update

vendor: Gopkg.lock
	@dep ensure

.PHONY: clean-vendor
clean-vendor:
	@rm -rf vendor

# Testing
.PHONY: test
test: vendor
	@go test -short

.PHONY: test-mysql
test-mysql: vendor
	@go test -v ./internal/mysql -scripts=file://migrations -db-user $(DB_USER) -db-password $(DB_PASSWORD) -db-address $(DB_ADDRESS) -db-name $(DB_NAME)

# Database Migration
.PHONY: migrate-prepare
migrate-prepare:
	@go get -u -d github.com/mattes/migrate/cli github.com/go-sql-driver/mysql
	@go build -tags 'mysql' -o /usr/local/bin/migrate github.com/mattes/migrate/cli

.PHONY: migrate-up
migrate-up:
	@migrate -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_ADDRESS))/$(DB_NAME)?multiStatements=true" -path=internal/mysql/migrations up

.PHONY: migrate-down
migrate-down:
	@migrate -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_ADDRESS))/$(DB_NAME)?multiStatements=true" -path=internal/mysql/migrations down

.PHONY: migrate-drop
migrate-drop:
	@migrate -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_ADDRESS))/$(DB_NAME)?multiStatements=true" -path=internal/mysql/migrations drop

# Upstream Services
.PHONY: docker-mysql-up
docker-mysql-up:
	@docker run --rm -d --name jobexec-mysql -p ${DB_PORT}:3306 -e MYSQL_DATABASE=$(DB_NAME) -e MYSQL_USER=$(DB_USER) -e MYSQL_PASSWORD=$(DB_PASSWORD) -e MYSQL_ROOT_PASSWORD=rootsecret mysql && docker logs -f jobexec-mysql

.PHONY: docker-mysql-down
docker-mysql-down:
	@docker stop jobexec-mysql
