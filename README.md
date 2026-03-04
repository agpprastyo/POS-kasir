# **POS Kasir (Point of Sales System)**

[![CI/CD](https://github.com/agpprastyo/POS-kasir/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/agpprastyo/POS-kasir/actions/workflows/ci-cd.yml)

## Overview

**POS Kasir** is a modern, high-performance Fullstack Point of Sales application designed to streamline retail operations. It provides a robust solution for managing products, processing orders, handling payments (including Digital Payments via Midtrans), and analyzing sales performance.

Built as a **single-port deployment** — the Go backend serves both the REST API and the React SPA frontend on port `8080`.

> **Note:** This project serves as a portfolio showcase demonstrating full-stack development capabilities, system architecture design, and integration of third-party services.

## ✨ Key Features

- **User Management & RBAC** — JWT authentication with role-based access control (Admin, Manager, Cashier)
- **Inventory Management** — Products, categories, variants/options, stock history
- **Order Processing** — Cart system, order workflow, operational status tracking
- **Digital Payments** — Integrated with Midtrans Payment Gateway (sandbox mode)
- **Shift Management** — Cashier shift tracking with cash transactions
- **Cloud Storage** — Cloudflare R2 / MinIO (S3-compatible) for product images
- **Dashboard & Analytics** — Sales reports, cashier performance, profit summaries
- **Activity Logging** — Complete audit trails
- **Multi-language** — i18n support (English / Indonesian)
- **Promotions** — Discount and promotion management
- **Thermal Printing** — ESC/POS receipt printing support

## Tech Stack

### Backend
- **Go 1.25** + [Fiber v3](https://gofiber.io/) (HTTP framework)
- **PostgreSQL** + [sqlc](https://sqlc.dev/) (type-safe SQL)
- **JWT** authentication + RBAC middleware
- **Swagger** auto-generated API docs

### Frontend
- **React 19** + [TanStack Router](https://tanstack.com/router) (SPA)
- **Vite 7** (build tool)
- **Tailwind CSS 4** + [shadcn/ui](https://ui.shadcn.com/) (UI components)
- **TanStack Query** (data fetching)
- **OpenAPI Generator** (typed API client)

### Infrastructure
- **Docker** multi-stage build (final image ~30MB Alpine)
- **GitHub Actions** CI/CD → Docker image to GHCR
- **PostgreSQL 15** + **MinIO** (S3-compatible storage)

## Architecture

```
┌──────────────────────────────┐
│         Go Backend :8080     │
│                              │
│  /api/v1/*    → REST API     │
│  /swagger/*   → API Docs     │
│  /healthz     → Health Check │
│  /*           → React SPA    │
│               (web/dist/)    │
└──────────┬───────────────────┘
           │
     ┌─────┴─────┐
     │ PostgreSQL │   MinIO (S3)
     └───────────┘
```

## 🚀 Quick Start (Docker)

Cara tercepat menjalankan seluruh aplikasi:

```bash
git clone https://github.com/agpprastyo/POS-kasir.git
cd POS-kasir

# 1. Copy dan edit environment
cp .env.example .env
# Edit .env → isi DB_PASSWORD dan JWT_SECRET

# 2. Jalankan seluruh stack
docker compose up -d

# 3. Buka aplikasi
# App     → http://localhost:8080
# Swagger → http://localhost:8080/swagger/index.html
# MinIO   → http://localhost:9001
```

### Env yang WAJIB diisi:

| Variable | Keterangan |
|----------|-----------|
| `DB_PASSWORD` | Password PostgreSQL |
| `JWT_SECRET` | Secret key untuk token (`openssl rand -hex 32`) |

Env lainnya sudah memiliki default yang aman. Lihat [`.env.example`](.env.example) untuk daftar lengkap.

## 🛠️ Development Setup

### Prerequisites

- **Go** 1.25+
- **Node.js** 22+ (npm)
- **PostgreSQL** 15+
- **Docker** & Docker Compose (untuk infra)

### 1. Setup Infrastructure

```bash
# Jalankan PostgreSQL + MinIO via Docker
docker compose -f docker-compose-infra.yaml up -d
```

### 2. Setup Environment

```bash
cp .env.example .env
# Edit .env:
#   DB_HOST=localhost
#   DB_PASSWORD=<password dari docker-compose-infra>
#   JWT_SECRET=<generate dengan: openssl rand -hex 32>
```

### 3. Run Backend

```bash
# Install air untuk hot-reload (opsional)
go install github.com/air-verse/air@latest

# Run dengan hot-reload
air

# Atau tanpa hot-reload
go run ./cmd/app
```

Backend berjalan di `http://localhost:8080`

### 4. Run Frontend (Dev Mode)

```bash
cd web
npm install --legacy-peer-deps
npm run dev
```

Frontend dev server di `http://localhost:5173` (dengan proxy ke backend 8080)

### 5. Build Frontend (SPA)

```bash
cd web
npm run build
# Output: web/dist/
```

Setelah build, akses `http://localhost:8080` — Go backend serve SPA langsung.

## Useful Commands

```bash
# Database
make migrate-up          # Jalankan migrations
make migrate-down-one    # Rollback 1 migration
make migrate-create name=add_xxx_table   # Buat migration baru
make seed                # Seed sample data

# Code Generation
make sqlc-generate       # Generate Go code dari SQL
make swag                # Generate Swagger docs + API client

# Docker
docker compose up -d              # Full stack (app + DB + MinIO)
docker compose -f docker-compose-infra.yaml up -d   # Infra saja
```

## CI/CD

Pipeline berjalan via **GitHub Actions**:

| Trigger | Job | Keterangan |
|---------|-----|-----------|
| Push/PR ke `master` | **test** | Go vet, Go test, FE build |
| Tag `v*.*.*` | **test** + **build-and-push** | Build Docker image → push ke GHCR |

### Release Flow

```bash
# 1. Tag versi baru
git tag -a v1.2.0 -m "v1.2.0: description"
git push origin v1.2.0

# 2. CI otomatis build & push ke:
#    ghcr.io/agpprastyo/pos-kasir:1.2.0
#    ghcr.io/agpprastyo/pos-kasir:1.2
#    ghcr.io/agpprastyo/pos-kasir:latest
```

## Project Structure

```
.
├── cmd/
│   ├── app/              # Main server entry point
│   └── seeder/           # Database seeder
├── config/               # Configuration loading
├── internal/             # Business logic (auth, orders, products, etc.)
├── pkg/                  # Shared libraries (logger, JWT, R2, Midtrans)
├── server/
│   ├── server.go         # App init, middleware, lifecycle
│   ├── routes.go         # API route registration
│   └── frontend.go       # SPA static file serving
├── sqlc/                 # SQL queries, schema, migrations
├── web/                  # Frontend (React SPA)
│   ├── src/routes/       # TanStack Router file-based routes
│   ├── src/lib/api/      # Generated API client
│   └── dist/             # Built SPA output (gitignored)
├── Dockerfile            # Multi-stage build (Node + Go → Alpine)
├── docker-compose.yml    # Full stack deployment
└── Makefile              # Command shortcuts
```

## 📸 Screenshots

| Login Page | Dashboard |
| :----: | :----: |
| ![Login Page](screenshots/01_login.png) | ![Dashboard](screenshots/02_dashboard.png) |

| Point of Sales (POS) | Payment & Checkout |
| :----: | :----: |
| ![POS](screenshots/03_pos.png) | ![Payment](screenshots/04_payment.png) |

| Transaction History | Product Management |
| :----: | :----: |
| ![Transaction History](screenshots/05_transaction.png) | ![Product Management](screenshots/06_product.png) |

| Reports & Analytics | Settings |
| :----: | :----: |
| ![Reports](screenshots/07_reports.png) | ![Settings](screenshots/08_settings.png) |

## API Documentation

Auto-generated Swagger documentation available at:

- **Local:** http://localhost:8080/swagger/index.html

## License

This project is licensed under the [MIT License](LICENSE).

## Author

**Agung Prasetyo**

- GitHub: https://github.com/agpprastyo
- LinkedIn: https://www.linkedin.com/in/agprastyo
- Portfolio: https://portfolio.agprastyo.me
