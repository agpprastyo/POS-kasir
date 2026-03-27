# 

|  | Agung Prasetyo Backend Developer Sleman, Yogyakarta 52883 | [prasetyo.agpr@gmail.com](mailto:prasetyo.agpr@gmail.com) LinkedIn: [linkedin.com/in/agprastyo](http://linkedin.com/in/agprastyo) | Portfolio: [http://portfolio.agprastyo.me](http://portfolio.agprastyo.me) GitHub: [github.com/agpprastyo](http://github.com/agpprastyo) |
| :---- | :---: |

## **Summary**

## ---

Result-oriented Backend Engineer specializing in **Go (Golang)**, **PostgreSQL**, and distributed systems architecture. Expert in designing high-performance REST APIs with **Fiber v3**, implementing **real-time synchronization via WebSockets**, and optimizing data retrieval using **Redis cache-aside strategies**. Proficient in full-stack development using **React 19** and **TanStack** ecosystem. Dedicated to modern engineering practices including **Structured Logging (slog)**, **Sentry observability**, and automated **CI/CD pipelines**. Passionate about building secure, scalable, and audit-ready systems that deliver premium user experiences.

## **Skills**

## ---

* **Backend Development:** **Go (Golang)**, Fiber v3, Clean Architecture, **Real-time WebSockets**, Structured Logging (slog), Sentry Observability, JWT & Refresh Tokens, RBAC (Role-Based Access Control)
* **Databases & Caching:** PostgreSQL, **sqlc (Type-safe SQL)**, Redis (Caching & Rate Limiting), Database Schema Design, Managed Migrations
* **Architecture & DevOps:** Single-port Deployments, **GitHub Actions CI/CD**, Docker & Docker Compose, **Tailscale SSH**, Swagger/OpenAPI, Git Workflow (v1.4.0+ semantic tagging)
* **Cloud & Storage:** **Cloudflare R2** / MinIO (S3-Compatible), Object Storage Lifecycle Management
* **Frontend (Strong Supporting):** **React 19**, TypeScript, **TanStack Router & Query**, Tailwind CSS 4, shadcn/ui, i18next

## **Projects**

## ---

**POS Kasir – Full Stack Point of Sales | Personal Project** | [Source Code](https://github.com/agpprastyo/POS-kasir) | [Live Demo](https://pos-kasir.agprastyo.me)  
*June 2025 – Present*  
**Tech Stack:** Go (Fiber v3), React 19, PostgreSQL (sqlc), Redis, WebSocket, Sentry, Cloudflare R2, Docker, GitHub Actions  
Engineered a production-ready POS system with a focus on real-time state synchronization, high-performance analytics, and observability.

*   Implemented **Real-time Synchronization** using a global WebSocket Hub to broadcast order updates and shift status across multiple cashier instances.
*   Optimized **Analytics Dashboard** performance by implementing **Redis cache-aside logic** for complex sales and product performance reports.
*   Enhanced system reliability with **Observability** upgrades, including `log/slog` structured logging and **Sentry SDK** for request-scoped error tracking and panic recovery.
*   Designed a secure **RBAC system** integrated into both Backend middleware and Frontend UI components (shadcn/ui), ensuring granular permission-based access.
*   Architected a **Single-Port Deployment** strategy where a high-performance Go binary serves both the REST API and the React SPA, simplifying infrastructure.
*   Developed an **Automated Maintenance Pipeline** featuring a daily cron-based database reset (Wipe & Re-seed) to ensure demo environment consistency.
*   Established a robust **CI/CD workflow** via GitHub Actions, including automated builds, Docker image versioning (v1.4.0+), and Tailscale-secured deployments.
*   Achieved **100% test pass rate** for core business logic using `mockgen` and `pgxmock` before major releases.

**KirimKarya \- Photo Delivery & Client Proofing Platform** | [Source Code](https://github.com/agpprastyo/KirimKarya)  
*Mar 2026 \- Present*   
**Tech Stack:** Bun, TypeScript, Hono, PostgreSQL (Drizzle ORM), Redis, BullMQ, S3, SvelteKit  
Engineered a high-performance backend architecture designed to support image delivery, background processing, and scalable file storage.

* Designed a **monorepo architecture** separating API services, frontend applications, and background worker services.  
* Built an **asynchronous image processing pipeline** using Redis and BullMQ to handle heavy tasks such as watermarking and thumbnail generation.  
* Implemented **background workers** to offload CPU-intensive image processing and maintain fast API response times.  
* Integrated **S3-compatible object storage** for scalable image uploads and file management.  
* Developed automated **cron-based lifecycle jobs** to remove expired gallery assets from storage and database.  
* Secured API endpoints using **Better Auth session management, RBAC authorization, and strict request validation using Zod**.

**Experience**

## ---

**Mobile Flutter Developer Intern**  
*PT. Solusi Digital Handal* — Yogyakarta  
*February 2024 – May 2024*

* Converted **20+ Figma UI designs** into pixel-perfect Flutter widgets across various screen sizes.  
* Integrated multiple REST APIs into mobile applications for real-time data retrieval.  
* Implemented **Bloc state management** to improve maintainability and code structure.

## **Additional Experience**

## ---

**Staff Photographer**  
*Laris Studio* — Sleman, Yogyakarta  
*July 2021 – Present*

* Handled client projects from planning to delivery, ensuring high customer satisfaction.  
* Managed time and priorities effectively to meet strict deadlines.  
* Developed strong discipline, responsibility, and client communication skills.

## **Education**

## ---

***Bachelor's Degree in Informatics (Fully Online / Distance Learning)***   
*Universitas Siber Muhammadiyah Mar 2023 – Present*