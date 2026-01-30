FROM oven/bun:1 AS builder

WORKDIR /app


COPY package.json bun.lock ./


RUN bun install --frozen-lockfile

COPY . .

RUN bun run build

FROM nginx:alpine

# Remove default nginx config
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built assets from builder stage
COPY --from=builder /app/dist/client /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
