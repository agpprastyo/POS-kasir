# --- Stage 1: Build Frontend (Svelte) ---
FROM node:20-alpine AS frontend
WORKDIR /app/web

# Salin file package dan install dependensi
COPY web/package.json web/package-lock.json ./
RUN npm install

# Salin sisa source code frontend
COPY web/ ./

# Bangun aset statis
# Asumsi output build ada di direktori 'build'
RUN npm run build

# --- Stage 2: Build Backend (Go) ---
FROM golang:1.25-alpine AS backend
WORKDIR /app

# Salin file mod dan unduh dependensi
COPY go.mod go.sum ./
RUN go mod download

# Salin sisa source code backend
COPY . .

# Bangun binary Go
# CGO_ENABLED=0 penting untuk static linking di Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/main ./cmd/app/main.go

# --- Stage 3: Final Production Image ---
FROM alpine:latest
WORKDIR /app

# Install certificate dan timezone
RUN apk --no-cache add tzdata ca-certificates

# Install 'dockerize' untuk menunggu Postgres & Minio siap
ENV DOCKERIZE_VERSION v0.7.0
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

# Salin migrasi (jika aplikasi Anda menjalankannya saat startup)
COPY --from=backend /app/sqlc/migrations ./sqlc/migrations

# Salin aset statis Svelte
# Server Go Anda harus dikonfigurasi untuk menyajikan file dari './public'
COPY --from=frontend /app/web/build ./public

# Salin binary Go
COPY --from=backend /app/main .

# Port yang diekspos oleh aplikasi Go
EXPOSE 8080

# Jalankan aplikasi
# Ini akan menunggu DB_HOST dan MINIO_ENDPOINT merespon sebelum menjalankan /app/main
# Pastikan variabel ini ada di .env Anda!
CMD ["dockerize", \
     "-wait", "tcp://${DB_HOST}:${DB_PORT}", \
     "-wait", "tcp://${MINIO_ENDPOINT}", \
     "-timeout", "60s", \
     "/app/main"]