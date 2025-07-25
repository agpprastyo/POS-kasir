# POS Kasir (Point of Sale System)

A modern, modular, and extensible Point of Sale (POS) system built with Go and Svelte. This project is designed for retail, restaurant, and small business environments, providing robust order management, inventory tracking, and payment processing.
<hr></hr>

## Features 
* User Authentication & Authorization 
* Product & Category Management 
* Order Management (Dine-in, Takeaway, Delivery)
* Inventory Management 
* Promotion & Discount Handling 
* Multiple Payment Methods (Cash, QRIS, etc.)
* Activity Logging & Audit Trail 
* Reporting (Sales, Inventory, Profit/Loss)
* RESTful API 
* Responsive Web Frontend (Svelte + Tailwind CSS)
* Dockerized for Easy Deployment
<hr></hr>

## Tech Stack
* Backend: Go (Golang), PostgreSQL, SQLC 
* Frontend: Svelte, Tailwind CSS, Vite 
* API: RESTful, JWT Auth 
* DevOps: Docker, Makefile

<hr></hr>

## Project Structure
```
cmd/            # Application entrypoints
config/         # Configuration and helpers
internal/       # Business logic (auth, orders, products, etc.)
mocks/          # Mock implementations for testing
pkg/            # Shared packages (database, logger, middleware, etc.)
repository/     # SQLC generated queries and models
server/         # HTTP server, routes, health checks
sqlc/           # SQL migrations and SQLC config
web/            # Frontend (Svelte)
```

<hr></hr>

## API Overview

### Authentication
* `POST /auth/login` — User login
* `POST /auth/register` — User registration

### Products & Categories
* `GET /products `— List products
* `POST /products` — Create product
* `GET /categories` — List categories

### Orders

* `POST /orders` — Create order
* `PATCH /orders/{id}/items` — Update order items
* `POST /orders/{id}/process-payment` — Process payment (QRIS/Online)
* `POST /orders/{id}/complete-manual-payment` — Complete manual payment (Cash)
* `PATCH /orders/{id}/status` — Update order status
* `POST /orders/{id}/cancel` — Cancel order

### Payments
* `GET /payment-methods` — List payment methods

### Promotions
* `GET /promotions` — List promotions

<hr></hr>

## Getting Started

### Prerequisites

* Docker & Docker Compose
* Go 1.21+
* Node.js 18+

### Quick Start