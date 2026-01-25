# ==========================================
# Stage 1: Builder
# ==========================================
FROM oven/bun:1 AS builder

WORKDIR /app

# Copy dependency files
COPY package.json bun.lock ./

# Install dependencies (frozen lockfile for reproducibility)
RUN bun install --frozen-lockfile

# Copy source code
COPY . .

# Build the application
# This will output static files to the 'dist' folder
RUN bun run build

# ==========================================
# Stage 2: Runtime (Nginx)
# ==========================================
FROM nginx:alpine

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html



EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
