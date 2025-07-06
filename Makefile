# Makefile for the VELORAS CLI tool

# Load .env file if it exists
ifneq (,$(wildcard .env))
	include .env
	export
endif

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

# 
GOOSE_DRIVER = postgres
BINARY_NAME = veloras-cli
MAIN_RUN = ./cmd/server/main.go
SWAG_DOCS = ./cmd/swag/docs
MIGRATIONS_DIR = ./internal/db/schemas
GOOSE_DBSTRING = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Default target is to build the binary
all: build

start:
	go run $(MAIN_RUN)

dev:
	air

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Cross-platform builds
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux-amd64 -v

build-mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-darwin-amd64 -v

# Build for all platforms
build-all: build-linux build-windows build-mac

# Install to GOPATH/bin
install:
	$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) -v

# Docker down
down:
	docker-compose down

# SQLC generator
sqlc:
	sqlc generate

# swagger
swag:
	@echo "Generating Swagger documentation..."
	swag init -g $(MAIN_RUN) -o $(SWAG_DOCS)


# migrations
create-migration:
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

migrate-up-one:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	goose -dir=$(MIGRATIONS_DIR) up-by-one

migrate-up:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	goose -dir=$(MIGRATIONS_DIR) up

migrate-down:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	goose -dir=$(MIGRATIONS_DIR) down

.PHONY: all build clean build-linux build-windows build-mac build-all install start