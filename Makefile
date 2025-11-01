# ==============================================================================
# Variabel
# ==============================================================================

# Variabel Migrasi
MIGRATE=migrate
MIGRATIONS_PATH=./sqlc/migrations

# Muat file .env dan ekspor variabel-variabelnya
include .env
export $(shell xargs < .env)

# Bangun DB_URL dari variabel .env
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Variabel Docker
COMPOSE_INFRA = docker-compose -f docker-compose.infra.yaml
COMPOSE_APP   = docker-compose -f docker-compose.app.yaml
COMPOSE_DEV   = docker-compose -f docker-compose.dev.yaml
DOCKER_NETWORK = kasir-net

# Menandai target yang bukan file
.PHONY: help all setup network infra-up infra-down prod prod-up prod-down prod-logs prod-nuke dev dev-d dev-down dev-logs sqlc-generate migrate-version migrate-up migrate-down migrate-down-one migrate-force migrate-create

# ==============================================================================
# Target Utama (Bantuan)
# ==============================================================================

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  setup         : (Internal) Memeriksa dan membuat jaringan Docker."
	@echo ""
	@echo "  infra-up      : 🚀 Memulai infrastruktur (Postgres & Minio) di background."
	@echo "  infra-down    : 🛑 Menghentikan infrastruktur (Postgres & Minio)."
	@echo ""
	@echo "  prod          : 🏭 Membangun & memulai aplikasi (PROD) & menampilkan log (blocking)."
	@echo "  prod-up       : 🏭 Membangun & memulai aplikasi (PROD) di background."
	@echo "  prod-down     : 🛑 Menghentikan aplikasi (PROD)."
	@echo "  prod-logs     : 📜 Menampilkan log aplikasi (PROD) (blocking)."
	@echo "  prod-nuke     : 💥 Menghentikan aplikasi (PROD) DAN infrastruktur."
	@echo ""
	@echo "  dev           : 🧑‍💻 Membangun & memulai lingkungan (DEV) di foreground (hot-reload)."
	@echo "  dev-d         : 🧑‍💻 Membangun & memulai lingkungan (DEV) di background."
	@echo "  dev-down      : 🛑 Menghentikan lingkungan (DEV)."
	@echo "  dev-logs      : 📜 Menampilkan log (DEV) (blocking)."
	@echo ""
	@echo "  sqlc-generate : Menjalankan 'sqlc generate'."
	@echo "  migrate-up    : Menjalankan semua migrasi 'up'."
	@echo "  migrate-down  : Menjalankan semua migrasi 'down'."
	@echo "  migrate-create name=<nama> : Membuat file migrasi baru."

# ==============================================================================
# Target Setup Docker
# ==============================================================================

# Target 'setup' hanya menjalankan 'network'
setup: network

# Memeriksa apakah jaringan ada, jika tidak, buat jaringan tersebut.
# @ di depan perintah menyembunyikan perintah itu sendiri dari output.
network:
	@docker network inspect $(DOCKER_NETWORK) >/dev/null 2>&1 || \
		(echo "Membuat jaringan Docker: $(DOCKER_NETWORK)" && docker network create $(DOCKER_NETWORK))

# ==============================================================================
# Target Infrastruktur (DB, Minio)
# ==============================================================================

infra-up: setup
	@echo "🚀 Memulai infrastruktur (Postgres & Minio)..."
	@$(COMPOSE_INFRA) up -d

infra-down:
	@echo "🛑 Menghentikan infrastruktur..."
	@$(COMPOSE_INFRA) down

# ==============================================================================
# Target Production
# ==============================================================================

# Target 'prod' akan memulai DAN menampilkan log
prod: prod-up prod-logs

prod-up: infra-up
	@echo "🏭 Membangun dan memulai aplikasi production..."
	@$(COMPOSE_APP) up --build -d

prod-down:
	@echo "🛑 Menghentikan aplikasi production..."
	@$(COMPOSE_APP) down

prod-logs:
	@echo "📜 Menampilkan log production (Ctrl+C untuk keluar)..."
	@$(COMPOSE_APP) logs -f

prod-nuke:
	@echo "💥 Menghentikan aplikasi production DAN infrastruktur..."
	@$(COMPOSE_APP) down
	@$(COMPOSE_INFRA) down

# ==============================================================================
# Target Development
# ==============================================================================

# 'make dev' akan berjalan di foreground, menampilkan semua log hot-reload
dev: infra-up
	@echo "🧑‍💻 Memulai lingkungan development (foreground, Ctrl+C untuk berhenti)..."
	@$(COMPOSE_DEV) up --build

# 'make dev-d' akan berjalan di background
dev-d: infra-up
	@echo "🧑‍💻 Memulai lingkungan development (detached/background)..."
	@$(COMPOSE_DEV) up --build -d

dev-down:
	@echo "🛑 Menghentikan lingkungan development..."
	@$(COMPOSE_DEV) down

dev-logs:
	@echo "📜 Menampilkan log development (Ctrl+C untuk keluar)..."
	@$(COMPOSE_DEV) logs -f

# ==============================================================================
# Target Database & SQLC (Milik Anda)
# ==============================================================================

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