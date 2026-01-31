# **🛒 POS Kasir (Point of Sales System)**

\<div align="center"\>

\<\!-- Live Demo Buttons \--\>

\<a href="https://pos-kasir.agprastyo.me/"\>

\<img src="https://www.google.com/search?q=https://img.shields.io/badge/🚀\_Live\_Frontend-Visit\_App-2ea44f?style=for-the-badge\&logo=vercel" alt="Live Frontend"\>

\</a\>

\<a href="https://api-pos.agprastyo.me/api/v1"\>

\<img src="https://www.google.com/search?q=https://img.shields.io/badge/⚙️\_Live\_API-Base\_URL-orange?style=for-the-badge\&logo=postman" alt="Live API"\>

\</a\>

\</div\>

## **📖 Overview**

**POS Kasir** is a modern, high-performance Fullstack Point of Sales application designed to streamline retail operations. It provides a robust solution for managing products, processing orders, handling payments (including Digital Payments via Midtrans), and analyzing sales performance.

Built with **scalability** and **type-safety** in mind, the backend leverages **Golang** with **Fiber** and **sqlc**, while the frontend offers a seamless user experience using the bleeding-edge **TanStack Start** framework powered by **Bun**.

**Note:** This project serves as a portfolio showcase demonstrating full-stack development capabilities, system architecture design, and integration of third-party services.

## **✨ Key Features**

### **🏢 Core Functionality**

* **User Management & RBAC:** Secure authentication with JWT. Role-based access control for Admins and Cashiers.
* **Inventory Management:** Create, update, and organize products with categories. Support for product variants/options.
* **Order Processing:** Efficient cart system and order placement workflow.
* **Transactions:** Detailed transaction history and receipt generation.

### **🚀 Advanced Features**

* **Digital Payments:** Integrated with **Midtrans Payment Gateway** for seamless cashless transactions.
* **Cloud Storage:** Integration with **Cloudflare R2** for efficient and scalable product image storage.
* **Dashboard & Analytics:** Comprehensive reports on sales, cashier performance, and popular products.
* **Activity Logging:** Complete audit trails for tracking system changes and user activities.
* **Multi-language Support:** Frontend i18n support (English/Indonesian).

## **🛠️ Tech Stack**

### **Backend (API)**

* **Language:** [Go (Golang)](https://www.google.com/search?q=https://go.dev/)
* **Framework:** [Fiber v2](https://www.google.com/search?q=https://gofiber.io/) \- High-performance web framework.
* **Database:** PostgreSQL.
* **ORM/Query Builder:** [sqlc](https://www.google.com/search?q=https://sqlc.dev/) \- For generating type-safe Go code from SQL.
* **Migrations:** Golang Migrate.
* **Docs:** Swagger (Swaggo).
* **Utils:** Viper (Config), Zap (Logging).

### **Frontend (Web)**

* **Runtime:** [Bun](https://www.google.com/search?q=https://bun.sh/)
* **Framework:** [TanStack Start](https://www.google.com/search?q=https://tanstack.com/start/latest) (React).
* **State & Data Fetching:** [TanStack Query](https://www.google.com/search?q=https://tanstack.com/query/latest).
* **UI Component:** [Shadcn UI](https://www.google.com/search?q=https://ui.shadcn.com/) \+ Tailwind CSS.
* **Form Handling:** React Hook Form \+ Zod.
* **API Client:** OpenAPI Generator (generated from Backend Swagger).

### **Infrastructure & Tools**

* **Containerization:** Docker & Docker Compose.
* **Hot Reload:** Air (Backend).
* **Automation:** Makefile.

## **📂 Project Structure**

.  
├── cmd/                \# Main applications entry points  
│   ├── app/            \# Main server application  
│   └── seeder/         \# Database seeder  
├── config/             \# Configuration loading logic  
├── internal/           \# Private application and business logic  
│   ├── auth/           \# Authentication logic  
│   ├── orders/         \# Order processing  
│   ├── products/       \# Product management  
│   ├── repository/     \# Generated sqlc code  
│   └── ...  
├── pkg/                \# Public library code (Logger, Midtrans, R2, Utils)  
├── sqlc/               \# SQL queries and schema  
├── web/                \# Frontend application (TanStack Start)  
├── docker-compose.yml  \# Docker orchestration  
└── Makefile            \# Command runner

## **🚀 Getting Started**

### **Prerequisites**

* **Go** 1.22+
* **Bun** 1.0+ (for frontend)
* **Docker** & **Docker Compose**
* **PostgreSQL** (if running locally without Docker)

### **1\. Clone the Repository**

git clone \[https://github.com/agpprastyo/POS-kasir.git\](https://github.com/agpprastyo/POS-kasir.git)  
cd POS-kasir

### **2\. Backend Setup**

**Using Docker (Recommended):**

The project includes a Makefile to simplify commands.

\# Start the database and backend services  
make up

\# Run database migrations  
make migrate-up

\# (Optional) Seed the database with dummy data  
make seed

**Manual Setup:**

1. Copy .env.example to .env and configure your Database, Midtrans, and Cloudflare R2 credentials.
2. Run go mod download.
3. Run the server: go run cmd/app/main.go.

### **3\. Frontend Setup**

Navigate to the web directory:

cd web

\# Install dependencies  
bun install

\# Setup Environment Variables  
cp .env.example .env  
\# Ensure VITE\_API\_BASE\_URL points to your backend (default: http://localhost:8080)

\# Run the development server  
bun dev

Open [http://localhost:3000](https://www.google.com/search?q=http://localhost:3000) to view the application.

## **📸 Screenshots**

\<\!-- Tip: Upload screenshots to your repo's 'assets' folder or an image host and link them here \--\>

| Dashboard | Order Page |
| :---- | :---- |
|  |  |

| Product Management | Mobile View |
| :---- | :---- |
|  |  |

## **🔌 API Documentation**

The backend includes auto-generated Swagger documentation.

* **Live Specs:** [https://api-pos.agprastyo.me/swagger/index.html](https://www.google.com/search?q=https://api-pos.agprastyo.me/swagger/index.html) *(Assuming swagger is accessible publicly)*
* **Live Base URL:** https://api-pos.agprastyo.me/api/v1

**Running Locally:**

Once the local server is running, visit:

http://localhost:8080/swagger/index.html

## **🤝 Contribution**

Contributions are welcome\! If you have suggestions or want to improve the codebase:

1. Fork the repository.
2. Create a feature branch (git checkout \-b feature/NewFeature).
3. Commit your changes.
4. Push to the branch.
5. Open a Pull Request.

## **📝 License**

This project is licensed under the [MIT License](https://www.google.com/search?q=LICENSE).