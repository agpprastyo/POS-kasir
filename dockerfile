# ==========================================
# Stage 1: Builder (Mengkompilasi Binary)
# ==========================================
# PENTING: Gunakan versi 1.24 sesuai go.mod kamu
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git (diperlukan jika ada dependency yang mengambil dari git langsung)
RUN apk add --no-cache git

# 1. Copy file dependency dulu (agar cache docker optimal)
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy seluruh source code
COPY . .

# 3. Build Binary
# - CGO_ENABLED=0: Wajib karena kita tidak pakai library C (pgx itu pure Go)
# - -ldflags="-s -w": Mengecilkan ukuran binary (menghapus debug symbol)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/app/main.go

# ==========================================
# Stage 2: Runtime (Menjalankan Aplikasi)
# ==========================================
FROM alpine:latest

WORKDIR /root/

# 1. Install Dependencies System
# - ca-certificates: WAJIB agar bisa connect HTTPS ke Supabase, Cloudflare R2, & Midtrans
# - tzdata: Agar jam di log/database sesuai (WIB)
RUN apk --no-cache add ca-certificates tzdata

# 2. Set Timezone ke Jakarta (Penting untuk POS/Transaksi)
ENV TZ=Asia/Jakarta

# 3. Copy Binary dari Stage Builder
COPY --from=builder /app/main .

# 4. Copy Folder Migrations (PENTING!)
# Karena kamu pakai golang-migrate, file SQL biasanya dibaca dari folder saat runtime
# Jika kode kamu pakai "embed", langkah ini bisa diskip. Jika baca file path, ini wajib.
COPY --from=builder /app/sqlc/migrations ./sqlc/migrations

# 5. Expose port & Run
EXPOSE 8080
CMD ["./main"]