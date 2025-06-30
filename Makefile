MIGRATE=migrate
MIGRATIONS_PATH=./sqlc/migrations

include .env
export $(shell xargs < .env)

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

migrate-version:
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(DB_URL) version

migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(DB_URL) up

migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(DB_URL) down

migrate-down-one:
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(DB_URL) down 1

migrate-force:
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(DB_URL) force $(version)


migrate-create:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

sqlc-generate:
	sqlc generate -f sqlc/sqlc.yaml