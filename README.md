# **POS Kasir (Point of Sales System)**

[![CI/CD](https://github.com/agpprastyo/POS-kasir/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/agpprastyo/POS-kasir/actions/workflows/ci-cd.yml)
![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-4169E1?logo=postgresql&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)

## Overview

**POS Kasir** is a modern, high-performance Fullstack Point of Sales application designed to streamline retail operations. It provides a robust solution for managing products, processing orders, handling payments (including Digital Payments via Midtrans), and analyzing sales performance.

Built as a **single-port deployment** — the Go backend serves both the REST API and the React SPA frontend on port `8080`.

> **Note:** This project serves as a portfolio showcase demonstrating full-stack development capabilities, system architecture design, and integration of third-party services.

## Live Deployment

- **App:** https://pos-kasir.agprastyo.me
- **Swagger docs:** https://pos-kasir.agprastyo.me/swagger/index.html

## ✨ Key Features

| Category | Feature |
|----------|---------|
| **Auth & Access** | JWT authentication, RBAC (Admin / Manager / Cashier), session management |
| **Inventory** | Products, categories, variants/options, stock history, image uploads, soft-delete & restore |
| **Orders** | Cart system, order workflow, operational status tracking, item updates |
| **Payments** | Manual cash/payment methods, Midtrans Payment Gateway (QRIS dynamic/static) |
| **Shift Management** | Cashier shift open/close, cash transactions, cash reconciliation |
| **Customers** | Customer registration and selection per order |
| **Promotions** | Percentage & fixed-amount discounts, scope (order/item), rules & targets |
| **Reports & Analytics** | Sales trends, product performance, cashier ranking, profit summary, payment distribution, low stock, shift summary |
| **Thermal Printing** | ESC/POS receipt printing, **auto network printer discovery** (TCP scan port 9100), Bluetooth support via Web Bluetooth |
| **Cloud Storage** | Cloudflare R2 / MinIO (S3-compatible) for product & variant images |
| **Activity Logging** | Complete audit trails with entity-level tracking |
| **Real-time Sync** | Global WebSocket Hub for instant cashier synchronization |
| **Redis Caching** | Cache-aside for optimized reporting performance |
| **Demo Maintenance**| Automated daily database reset (Wipe & Seed) at 01:00 AM |
| **Multi-language** | i18n support (English / Indonesian) with `react-i18next` |
| **Theming** | Light / Dark / System mode |

## Tech Stack

### Backend

| Technology | Purpose |
|-----------|---------|
| **Go 1.25** + [Fiber v3](https://gofiber.io/) | HTTP framework |
| **PostgreSQL 15** + [sqlc](https://sqlc.dev/) | Type-safe SQL code generation |
| **Redis** | Rate limiting, shift caching, report performance (cache-aside) |
| **JWT** + RBAC middleware | Authentication & authorization |
| **WebSocket** | Real-time state synchronization |
| **Sentry** + `slog` | Structured logging & error tracking |
| **Swagger** (swaggo) | Auto-generated API documentation |
| **ESC/POS** (`pkg/escpos`) | Raw TCP thermal receipt printing |

### Frontend

| Technology | Purpose |
|-----------|---------|
| **React 19** + [TanStack Router](https://tanstack.com/router) | SPA with file-based routing |
| **Vite 7** + **Bun** | Build tool & package manager |
| **Tailwind CSS 4** + [shadcn/ui](https://ui.shadcn.com/) | UI component library |
| **TanStack Query** | Data fetching & caching |
| **OpenAPI Generator** | Typed API client from Swagger spec |
| **react-i18next** | Internationalization (ID/EN) |
| **Recharts** | Dashboard charts & analytics |
| **Web Bluetooth API** | Direct Bluetooth printer connectivity |

### Infrastructure

| Technology | Purpose |
|-----------|---------|
| **Docker** multi-stage build | Final image ~30MB Alpine |
| **GitHub Actions** CI/CD | Build → Push to GHCR → Deploy via Tailscale SSH |
| **MinIO** | S3-compatible object storage |
| **Redis** | Caching, rate limiting, and session state |
| **Sentry** | Production error & panic monitoring |

## Architecture

```
┌─────────────────────────────────────────────────┐
│              Go Backend :8080                    │
│                                                 │
│  /api/v1/*       → REST API                     │
│  /swagger/*      → API Docs                     │
│  /healthz        → Health Check                 │
│  /*              → React SPA (web/dist/)        │
└────────┬────────────────┬───────────────────────┘
         │                │
    ┌────┴────┐     ┌─────┴─────┐    ┌───────────┐
    │ Postgres │     │   Redis   │    │ MinIO/R2  │
    └─────────┘     └───────────┘    └───────────┘
```

```
Frontend (React SPA)
├── Network Printer  →  Backend TCP Socket  →  Printer :9100
└── Bluetooth Printer  →  Web Bluetooth API  →  BLE Printer
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
- **Bun** (or Node.js 22+)
- **PostgreSQL** 15+
- **Redis**
- **Docker** & Docker Compose (untuk infra)

### 1. Setup Infrastructure

```bash
# Jalankan PostgreSQL + MinIO + Redis via Docker
docker compose -f docker-compose.infra.yml up -d

# Atau menggunakan Makefile
make dev-infra
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
bun install
bun run dev
```

Frontend dev server di `http://localhost:5173` (dengan proxy ke backend 8080)

### 5. Build Frontend (SPA)

```bash
cd web
bun run build
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

Pipeline berjalan via **GitHub Actions** dan dipicu dari setiap push/PR ke `master` atau tag semver `v*.*.*`.

| Trigger | Job | Keterangan |
|---------|-----|------------|
| Push/PR ke `master` | **test** | `go vet` dan build frontend (`bun run build` di `./web`) |
| Tag `v*.*.*` | **test** + **build-and-push** | Build dan push image ke `ghcr.io/agpprastyo/pos-kasir` dengan tag `X.Y.Z`, `X.Y`, dan `latest` |
| Tag `v*.*.*` | **deploy** | Koneksi ke VM lewat **Tailscale**, lalu `docker compose pull` + `docker compose up -d` di `/home/ubuntu` |

CI mengandalkan variabel berikut:

- `REGISTRY` (default `ghcr.io`) dan `IMAGE_NAME` (nama repository) untuk metadata Docker
- `TAILSCALE_AUTHKEY` untuk autentikasi Tailscale
- `VM_TAILSCALE_IP` dan `VM_SSH_PRIVATE_KEY` agar `appleboy/ssh-action` bisa mengakses host

## Project Structure

```
.
├── cmd/
│   ├── app/              # Main server entry point
│   └── seeder/           # Database seeder
├── config/               # Configuration loading
├── internal/             # Business logic modules
│   ├── activitylog/      # Audit trail logging
│   ├── categories/       # Product category CRUD
│   ├── cancellation_reasons/ # Order cancellation reasons
│   ├── common/           # Shared middleware (auth, RBAC, rate limit, idempotency)
│   ├── customers/        # Customer management
│   ├── orders/           # Order processing & workflow
│   ├── payment_methods/  # Payment method CRUD
│   ├── printer/          # Thermal printing + network discovery
│   ├── products/         # Product CRUD, variants, stock, images
│   ├── promotions/       # Promotion engine (rules, targets, discounts)
│   ├── report/           # Sales, profit, performance analytics
│   ├── settings/         # App settings (branding, printer config)
│   ├── shift/            # Cashier shift management + cash reconcile
│   ├── user/             # User CRUD, auth, JWT, avatar
│   └── websocket/        # Real-time synchronization hub
├── pkg/                  # Shared packages
│   ├── cache/            # Redis cache abstraction
│   ├── cloudflare-r2/    # S3-compatible object storage
│   ├── database/         # PostgreSQL connection + migration runner
│   ├── escpos/           # ESC/POS printer protocol
│   ├── logger/           # Structured logging
│   ├── payment/          # Midtrans payment gateway
│   ├── utils/            # JWT manager, helpers
│   └── validator/        # Request validation
├── server/
│   ├── server.go         # App init, DI container, lifecycle
│   ├── routes.go         # API route registration
│   └── frontend.go       # SPA static file serving
├── sqlc/                 # SQL queries, schema, migrations
├── web/                  # Frontend (React SPA)
│   └── src/
│       ├── components/   # Feature-based component architecture
│       │   ├── account/        # Profile & security
│       │   ├── activity-logs/  # Activity log viewer
│       │   ├── auth/           # Login form
│       │   ├── common/         # Shared/reusable components
│       │   ├── customers/      # Customer management
│       │   ├── dashboard/      # Dashboard widgets & charts
│       │   ├── orders/         # POS cart, product search, variants
│       │   ├── payment/        # Payment dialogs
│       │   ├── products/       # Product table, grid, form, filters
│       │   ├── promotions/     # Promotion management
│       │   ├── reports/        # Analytics charts & tables
│       │   ├── settings/       # Branding, printer, categories
│       │   ├── transactions/   # Transaction history & actions
│       │   ├── users/          # User management
│       │   └── ui/             # shadcn/ui primitives
│       ├── lib/
│       │   ├── api/            # Generated API client + query hooks
│       │   ├── auth/           # RBAC utilities
│       │   ├── locales/        # i18n translations (id.json, en.json)
│       │   └── printer/        # Frontend PrinterService (Bluetooth)
│       └── routes/             # TanStack Router file-based routes
├── Dockerfile            # Multi-stage build (Bun + Go → Alpine)
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
- **Production:** https://pos-kasir.agprastyo.me/swagger/index.html

## License

This project is licensed under the [MIT License](LICENSE).

## Author

**Agung Prasetyo**

- GitHub: https://github.com/agpprastyo
- LinkedIn: https://www.linkedin.com/in/agprastyo
- Portfolio: https://portfolio.agprastyo.me
