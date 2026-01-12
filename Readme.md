# TiRush

> **A High-Concurrency Event Booking System built with Go.** > Handling 100k+ concurrent requests with zero double-bookings.\_

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![Postgres](https://img.shields.io/badge/PostgreSQL-15-336791?style=flat&logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?style=flat&logo=redis)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat&logo=docker)

## 🚀 The Challenge

Building a ticketing system is easy. Building one that survives a **Flash Sale** is hard.
**TiRush** is a backend engine designed to handle massive spikes in traffic where thousands of users compete for the same seat simultaneously.

### Key Problems Solved:

- **Race Conditions:** Prevents double-booking using **Optimistic Locking** (Postgres) and Distributed Locks (Redis).
- **Thundering Herd:** Protects the database using **Request Coalescing (Singleflight)** and tiered caching.
- **Data Integrity:** Uses **Idempotency Keys** to ensure payments are processed exactly once, even during network failures.
- **Real-Time Updates:** Websockets push seat availability to clients instantly.

---

## 🏗️ Architecture (Modular Monolith)

- **Language:** Golang (Standard Library + Chi Router)
- **Database:** PostgreSQL (with `pgx` connection pooling)
- **Caching/Locking:** Redis
- **Containerization:** Docker & Docker Compose

### DB Schema Design

Yet to be implemented
