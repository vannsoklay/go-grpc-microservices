# ==============================
# MIGRATION CONFIG
# ==============================

DATABASE_URL=$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

MIGRATIONS_DIR=migrations
MIGRATE=migrate

# ==============================
# MIGRATION COMMANDS
# ==============================

.PHONY: migrate-create migrate-up migrate-down migrate-force migrate-version migrate-drop

## Create new migration
migrate-create:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

## Run all migrations
migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

## Rollback last migration
migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

## Force migration version (use carefully)
migrate-force:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(version)

## Show current migration version
migrate-version:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

## Drop everything (DANGEROUS)
migrate-drop:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" drop
