# Makefile for running golang-migrate commands
include .envrc
# Directory containing migration files
MIGRATE_PATH=./cmd/migrate/migrations

# Database connection string
# DATABASE_URL=postgres://admin:adminpassword@localhost/social?sslmode=disable

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make migrate-up       - apply all up migrations"
	@echo "  make migrate-down     - rollback last migration"
	@echo "  make migrate-version  - show current migration version"
	@echo "  make migrate-create NAME=xxx - create new migration files"

# Apply all pending migrations
.PHONY: migrate-up
migrate-up:
	migrate -path=$(MIGRATE_PATH) -database=$(DB_ADDR) up

# Rollback the last migration
.PHONY: migrate-down
migrate-down:
	migrate -path=$(MIGRATE_PATH) -database=$(DB_ADDR) down 1

# Show current migration version
.PHONY: migrate-version
migrate-version:
	migrate -path=$(MIGRATE_PATH) -database=$(DB_ADDR) version

# Create a new migration
# Usage: make migrate-create NAME=create_users
.PHONY: migrate-create
migrate-create:
ifndef NAME
	$(error NAME is not set. Usage: make migrate-create NAME=create_users)
endif
	migrate create -seq -ext sql -dir=$(MIGRATE_PATH) $(NAME)
