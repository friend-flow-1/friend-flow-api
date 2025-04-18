# Friend Flow API

A Go RESTful API built with Gin Gonic, PostgreSQL, Casbin RBAC, and JWT authentication. This project supports modules for user management and authentication.

---

## 🚀 Features

- JWT-based authentication
- Role-based access control (RBAC) with Casbin
- PostgreSQL as the main database
- ScyllaDB support for Casbin policies (optional)
- Modular structure (`auth`, `user`)
- Gin Gonic framework
- Environment-based config

---

## 🐳 Getting Started with Docker (PostgreSQL)

### 🔧 Run PostgreSQL using Docker

```bash
docker run --name friendflow-postgres -e POSTGRES_USER=haxxu -e POSTGRES_PASSWORD=User123 -e POSTGRES_DB=friendflow -p 5432:5432 -d postgres:latest
```

### 🔧 Run ScyllaDB using Docker

```bash
docker run --name scylla -p 9042:9042 -d scylladb/scylla:latest
```
