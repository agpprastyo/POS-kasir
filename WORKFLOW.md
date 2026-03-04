# Workflow Guide — POS Kasir

## 1. Development Workflow

### Branching Strategy (Git Flow)

```
master ─────────────────────────────────────────► (production-ready)
  │                                        ▲
  └──► develop ──────────────────────────── merge ──►
         │          ▲        ▲        ▲
         ├── feature/xxx ────┘        │
         ├── feature/yyy ─────────────┘
         └── bugfix/zzz ──────────────┘

master ◄── hotfix/critical ──► (langsung dari master)
```

| Branch | Sumber | Merge ke | Kapan |
|--------|--------|----------|-------|
| `master` | — | — | Selalu production-ready |
| `develop` | `master` | `master` | Kumpulan fitur siap release |
| `feature/*` | `develop` | `develop` | Fitur baru |
| `bugfix/*` | `develop` | `develop` | Perbaikan non-urgent |
| `hotfix/*` | `master` | `master` + `develop` | Perbaikan kritis di production |

---

### Membuat Fitur Baru

```bash
# 1. Buat branch dari develop
git checkout develop
git pull origin develop
git checkout -b feature/nama-fitur

# 2. Develop & commit
git add .
git commit -m "feat: deskripsi singkat"

# 3. Push & buat PR ke develop
git push origin feature/nama-fitur
# Buat Pull Request: feature/nama-fitur → develop

# 4. Setelah PR approved & merged, hapus branch
git checkout develop
git pull origin develop
git branch -d feature/nama-fitur
```

**Contoh commit messages:**
```
feat: add product barcode scanning
feat(orders): implement split payment
fix: resolve stock deduction race condition
refactor: extract payment logic to service layer
```

### Bugfix (Non-urgent)

```bash
git checkout develop
git checkout -b bugfix/deskripsi-bug

# Fix → commit → PR ke develop
git commit -m "fix: deskripsi perbaikan"
git push origin bugfix/deskripsi-bug
```

### Hotfix (Production Critical)

```bash
# 1. Buat dari master
git checkout master
git checkout -b hotfix/deskripsi-kritis

# 2. Fix → commit
git commit -m "fix: critical issue description"

# 3. Merge ke master + tag
git checkout master
git merge hotfix/deskripsi-kritis
git tag -a v1.2.1 -m "v1.2.1: hotfix description"
git push origin master --tags

# 4. Merge balik ke develop
git checkout develop
git merge hotfix/deskripsi-kritis
git push origin develop

# 5. Cleanup
git branch -d hotfix/deskripsi-kritis
```

### Release

```bash
# 1. Merge develop ke master
git checkout master
git merge develop --no-ff -m "Merge develop: release v1.3.0"

# 2. Tag & push → CI auto build Docker image
git tag -a v1.3.0 -m "v1.3.0: release description"
git push origin master --tags

# 3. Verifikasi
# → https://github.com/agpprastyo/POS-kasir/actions
# → ghcr.io/agpprastyo/pos-kasir:1.3.0
# → ghcr.io/agpprastyo/pos-kasir:latest
```

---

## 2. Production Setup

### Minimal VPS Requirements
- **OS:** Ubuntu 22.04+ / Debian 12+
- **RAM:** 1 GB minimum
- **Docker:** 24.0+
- **Port:** 8080 (atau behind reverse proxy)

### Step-by-step

```bash
# 1. Install Docker (jika belum)
curl -fsSL https://get.docker.com | sh

# 2. Clone repo (untuk docker-compose.yml dan .env.example)
git clone https://github.com/agpprastyo/POS-kasir.git
cd POS-kasir

# 3. Setup environment
cp .env.example .env
nano .env
```

**Edit `.env` — yang WAJIB diganti:**
```env
APP_ENV=production
DB_PASSWORD=<password-kuat>
JWT_SECRET=<openssl rand -hex 32>
```

**Opsional (jika pakai Midtrans / R2):**
```env
MIDTRANS_SERVER_KEY=<server-key>
MIDTRANS_IS_PROD=true
R2_ACCOUNT_ID=<account-id>
R2_ACCESS_KEY=<access-key>
R2_SECRET_KEY=<secret-key>
R2_PUBLIC_DOMAIN=https://storage.yourdomain.com
```

```bash
# 4. Jalankan
docker compose up -d

# 5. Verifikasi
docker compose ps          # Semua service "healthy"
curl http://localhost:8080/healthz   # {"status":"ok"}

# 6. (Opsional) Seed data awal
docker compose exec app ./pos-server seed
```

### Dengan Reverse Proxy (Nginx/Caddy)

```nginx
# /etc/nginx/sites-available/pos-kasir
server {
    listen 80;
    server_name pos.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Update Production

```bash
cd POS-kasir
docker compose pull        # Pull image terbaru
docker compose up -d       # Recreate container
docker image prune -f      # Hapus image lama
```

### Backup Database

```bash
# Backup
docker compose exec postgres pg_dump -U pos_user pos_kasir > backup_$(date +%Y%m%d).sql

# Restore
cat backup_20260304.sql | docker compose exec -T postgres psql -U pos_user pos_kasir
```

---

## 3. Portfolio Presentation (Interview)

### Talking Points

#### Architecture & Design Decisions

> "Saya membangun POS ini sebagai **monolith yang bersih** — Go backend melayani REST API dan React SPA dari satu port. Ini menyederhanakan deployment tanpa mengorbankan separation of concerns di kode."

- **Kenapa Go + Fiber?** → Performance tinggi, single binary, low memory (~20MB RSS)
- **Kenapa SPA embedded?** → Single port, zero CORS issues, Docker image < 30MB
- **Kenapa sqlc bukan ORM?** → Type-safe tanpa magic, SQL murni, compile-time checking

#### Technical Highlights

| Topik | Yang Bisa Dibahas |
|-------|-------------------|
| **Auth** | JWT + refresh token rotation, RBAC middleware (3 roles) |
| **API Design** | RESTful, OpenAPI spec, auto-generated client |
| **Database** | PostgreSQL, sqlc (type-safe), auto-migration |
| **Payment** | Midtrans integration, webhook handling |
| **File Upload** | S3-compatible (Cloudflare R2 / MinIO), presigned URLs |
| **CI/CD** | GitHub Actions, multi-stage Docker, GHCR |
| **Frontend** | React 19, TanStack Router (file-based), shadcn/ui, i18n |
| **Testing** | Go unit tests, mock interfaces, testcontainers |
| **Shift System** | In-memory cache untuk active shift, cash transactions |

#### Demo Flow (5-10 menit)

```
1. Buka app → Login sebagai Admin
2. Dashboard → Tunjukkan analytics & charts
3. Produk → CRUD + upload gambar + variants
4. POS → Buat order → Pilih produk → Checkout
5. Payment → Demo Midtrans sandbox payment
6. Reports → Sales report, cashier performance
7. Settings → Branding, printer config
8. Swagger → Tunjukkan API docs
9. Terminal → `docker compose up -d` ← single command deploy
```

#### Pertanyaan Interview yang Mungkin Muncul

| Pertanyaan | Jawaban Singkat |
|------------|-----------------|
| "Kenapa tidak pakai microservices?" | Portofolio project, monolith lebih tepat. Tapi bisa di-split berkat clean architecture |
| "Bagaimana handle concurrent orders?" | Database transactions + row-level locking di PostgreSQL |
| "Kenapa tidak pakai SSR?" | SPA cukup untuk POS (internal tool), SSR overkill |
| "Bagaimana keamanan JWT?" | Access token (24h) + refresh token rotation + httpOnly cookies |
| "Bagaimana deploy?" | `git tag → CI build Docker → push GHCR → docker compose pull` |

### Quick Demo Setup (Interview Hari H)

```bash
# Siapkan sebelum interview — pastikan semuanya jalan
cd POS-kasir
docker compose up -d
make seed                 # Data sample
open http://localhost:8080

# Credentials:
# Admin:   admin@pos.com / admin123
# Cashier: cashier@pos.com / cashier123
```

> **Tip:** Siapkan tab browser yang sudah terbuka: App, Swagger, GitHub repo, GitHub Actions.
