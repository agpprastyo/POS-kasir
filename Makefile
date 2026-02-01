MIGRATE=migrate
MIGRATIONS_PATH=./sqlc/migrations


include .env



DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: help all setup network infra-up infra-down prod prod-up prod-down prod-logs prod-nuke dev dev-d dev-down dev-logs sqlc-generate migrate-version migrate-up migrate-down migrate-down-one migrate-force migrate-create

# ==============================================================================
# Target Utama
# ==============================================================================

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "  sqlc-generate : Menjalankan 'sqlc generate'."
	@echo "  migrate-up    : Menjalankan semua migrasi 'up'."
	@echo "  migrate-down  : Menjalankan semua migrasi 'down'."
	@echo "  migrate-create name=<nama> : Membuat file migrasi baru."
	@echo "  seed          : Menjalankan seeder."
	@echo "  swag          : Menjalankan swag."


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
	@if [ -z "$(name)" ]; then \
		echo "Error: 'name' belum diatur. Contoh: make migrate-create name=add_user_table"; \
		exit 1; \
	fi
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

sqlc-generate:
	sqlc generate -f sqlc/sqlc.yaml

seed:
	go run ./cmd/seeder/main.go

swag:
	@echo "Generating Swagger docs to ./docs and web/apid-docs..."
	@swag init -g ./cmd/app/main.go -o ./docs --parseDependency --parseInternal
	@swag init -g ./cmd/app/main.go -o web/api-docs --parseDependency --parseInternal --outputTypes json
	@cd web && npm run api:gen

docker-be:
	docker compose -f docker-compose.backend.yml up -d --build backend