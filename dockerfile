FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Jakarta

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]