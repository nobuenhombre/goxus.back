# AUDIT REPORT: goxus backend (Go 1.26.1)

**Date:** 2025-06-03  
**Auditor:** Hermes Agent  
**Scope:** `/home/bookworker06JAN1979/Sources/golang.app/goxus/back`  
**Command:** `go test ./... -coverprofile=c.out`, `go vet ./...`, manual code review  
**Policy:** Code was reviewed but **not modified** during this audit.

---

## Table of Contents

1. [Overall Architecture](#1-overall-architecture)
2. [Dependency Injection & Wire](#2-dependency-injection--wire)
3. [Domain Layer (user)](#3-domain-layer-user)
4. [RBAC Service](#4-rbac-service)
5. [API Layer (Gin)](#5-api-layer-gin)
6. [Security](#6-security)
7. [Database & Migrations](#7-database--migrations)
8. [Testing](#8-testing)
9. [Dependencies](#9-dependencies)
10. [Code Style & Quality](#10-code-style--quality)
11. [Potential Bugs](#11-potential-bugs)
12. [Recommendations (Priority)](#12-recommendations-priority)
13. [Final Score](#13-final-score)

---

## 1. Overall Architecture

The project follows a classic layered architecture with clear separation of concerns:

```
cmd/goxus/           — entrypoint (main, Wire injector, app orchestrator)
internal/app/goxus/  — business logic & infrastructure
  cli/               — CLI flags
  config/            — YAML configuration
  domain/            — domain layer (user + orchestrator)
  api/server/        — HTTP server (Gin)
  cron-job/          — job scheduler
  log/               — log file management
internal/pkg/        — reusable packages
  db/goxus/          — xo-generated PostgreSQL repositories
  services/rbac/     — RBAC service
  hash/              — hashing (MD5)
pkg/tests/postgres/  — testcontainers helper
```

### Strengths

- Clear domain / API / infrastructure separation.
- Use of `internal/` to prevent external import.
- Single orchestrator (`App`) with 3 run modes: `init`, `service`, `default`.
- Graceful shutdown via `SIGINT`/`SIGTERM` with a 5-second timeout.
- All providers return `(Service, func(), error)` — Wire cleanup chain.

### Notes

- `App.Run()` uses a type assertion `a.cliConfig.(*cli.Config)` — safe because Wire builds it, but slightly brittle if the interface changes.

---

## 2. Dependency Injection & Wire

Google Wire v0.7.0 is used correctly:

- Every package exports a `ProviderSet` in `provider.go`.
- The injector in `wire.go` aggregates all sets — no logic, only `Build()`.
- `wire_gen.go` is auto-generated; cleanup functions are called in reverse construction order.

### Issue

- `ProvideAPI` in `server/provider.go` does:
  ```go
  new(configApp.Get().Hosts.API)
  ```
  This takes the address of a field inside the config. If the config is reused and mutated, a pointer alias could lead to a data race. Prefer copying the value.

---

## 3. Domain Layer (user)

### 3.1 Authorization Decorator

- `impl.go` — pure business logic, no auth.
- `authorized_service.go` — RBAC decorator implementing the same `Service` interface.
- `provider.go` assembles the chain: `NewAuthorized(New(...), rbac)`.

This is a correct and flexible pattern. Drawback:
- Every decorator method duplicates the template: `ActorIDFromContext` → `Check` → delegate. As the method count grows, a generic helper should be extracted.

### 3.2 Service Interface

- `Login`/`Logout`/`ValidateToken` bypass RBAC — correct.
- `GetByID` allows self-read (`actorID == id`) without permission — sensible for a SaaS admin panel.

### 3.3 Business Logic (impl.go)

- Soft-delete via `pq.NullTime` — correct.
- Duplicate email check on `Create` and `Update` — present.
- Self-deletion guard (`ErrCannotDeleteSelf`) — present.
- **Passwords are stored in plaintext** — see Security section.

---

## 4. RBAC Service

### 4.1 Data Model

Four tables:
- `rbac_roles`
- `rbac_permissions`
- `rbac_role_permissions` (many-to-many)
- `rbac_user_roles` (many-to-many)

### 4.2 Implementation

- Full CRUD for roles and permissions.
- `AssignPermissionsToRole` is idempotent (skips duplicates).
- `AssignRoleToUser` is strict (returns `ErrAlreadyExists` on duplicate).
- `DeleteRole` cascades `role_permissions` but blocks if role is assigned to any user.
- `DeletePermission` blocks if permission is assigned to any role.

### 4.3 Performance

- `CheckUserPermission`, `GetUserRoles`, `GetRolePermissions` call `GetAll()` on join tables and filter in-memory.
- At scale, this becomes a bottleneck. SQL queries with `JOIN` + `WHERE` should replace full table scans.
- `DeleteRole`/`DeletePermission` also scan entire join tables to check usage. `EXISTS` queries are preferable.

---

## 5. API Layer (Gin)

### 5.1 Routing

- `/health` — unversioned.
- `/api/v1/` — versioned group.
- Public and protected routes separated by `authMiddleware`.

### 5.2 Handlers

- DTO structs with Gin binding tags (`required`, `email`, `min`).
- `userToResponse` hides `password` — correct.
- `authContext` extracts the actor from `gin.Context` and places it into `context.Context` — correct; the domain does not depend on Gin.

### 5.3 Middlewares

- `CORSMiddleware` — `Allow-Origin: *` (acceptable for dev, risky for production).
- `APILoggerMiddleware` — logs method, path, status, latency.
- `AuthTokenMiddleware` — parses Bearer token, validates it, sets `user` and `token` in Gin context.

### 5.4 Response Statuses

- Login: `404` (not found), `403` (access denied), `503` (other).
- Logout: `401` if no token in context (unexpected after middleware), `404` if token not found (unlikely after middleware). `204 No Content` would be more idiomatic.
- CRUD: `201 Created`, `200 OK`, `404 Not Found`, `409 Conflict`. Delete returns `200` with a message instead of `204`.

### 5.5 API Gaps

- No pagination on `ListUsers`.
- No filtering or sorting.
- No rate limiting.
- No request ID / correlation ID in logs.
- No structured logging (JSON) — only `log.Printf`.

---

## 6. Security

### Critical

| # | Issue | Location | Details |
|---|-------|----------|---------|
| 1 | **Passwords stored in plaintext** | `impl.go:55`, `impl.go:219` | `Password: password // TODO: hash password before storing`; comparison is direct string equality. MD5 exists in `pkg/hash/md5.go` but is not used and is unsuitable for passwords. Use `bcrypt` or `argon2`. |
| 2 | **CORS Allow-Origin: *** | `middlewares.cors.go:8` | Acceptable for development; production must use a whitelist. |
| 3 | **No HTTPS / TLS** | `server.go:43` | `ListenAndServe()` instead of `ListenAndServeTLS`. Production requires TLS termination (reverse proxy or built-in). |

### High

| # | Issue | Details |
|---|-------|---------|
| 4 | **No rate limiting on login** | Brute-force possible. Add a per-IP or per-email rate limiter. |
| 5 | **Stateful tokens (UUID v4) in DB** | Stored in `users_tokens`, soft-deleted on logout. Fine for low-load admin panels, but `expires_at` and a cleanup job are needed. |

### Medium

| # | Issue | Details |
|---|-------|---------|
| 6 | **No input sanitization** | No XSS protection on `name`/`email` (lower risk with JSON API). No max length enforcement beyond Gin binding. |
| 7 | **No audit log** | Who created/deleted a user is not logged. |
| 8 | **`log.Fatal` in graceful shutdown** | `server.graceful.shutdown.go:27` — aborts the Wire cleanup chain. Use `log.Printf` and return the error. |

---

## 7. Database & Migrations

- 6 migrations via `golang-migrate`.
- xo-generated type-safe repositories.
- Soft-delete is application-level (`pq.NullTime`), not database-level.
- **Missing** partial unique index on `users(email) WHERE deleted_at IS NULL`.
  - `GetUserByEmail` can find a soft-deleted user.
  - Restoring a user may violate email uniqueness.
- `users_tokens` has no `expires_at` or TTL. Tokens live forever.

---

## 8. Testing

### Coverage (`go test ./... -cover`)

| Package | Coverage |
|---------|----------|
| `domain/user` | **86.4%** |
| `pkg/services/rbac` | **85.2%** |
| `config` | **41.4%** |
| Everything else | **0.0%** |
| **Total** | **19.0%** |

### Test Infrastructure

- Testcontainers (PostgreSQL) via `testpostgres.StartPostgresContainer`.
- Migrations applied in `TestMain`.
- `truncateAll` between tests — full table cleanup.
- `testify` assert/require.

### Structure

- `service_test.go` — tests for `impl` (raw).
- `authorized_service_test.go` — tests for the decorator (auth).
- No table-driven tests; each scenario is a separate function (readable but verbose).

### Gaps

- No tests for HTTP handlers (`0%` coverage).
- No tests for middlewares.
- No tests for graceful shutdown.
- No benchmark tests.
- No integration tests for Wire initialization.

---

## 9. Dependencies

### Direct

- `gin-gonic/gin` v1.12.0
- `google/wire` v0.7.0
- `robfig/cron/v3` v3.0.1
- `yaml.v3` v3.0.1
- `stretchr/testify` v1.11.1
- `testcontainers/testcontainers-go` v0.42.0
- `nobuenhombre/suikat` v0.0.170 (internal framework)

### Indirect

- `jackc/pgx/v5` v5.4.3 (via suikat)
- `lib/pq` v1.10.2 (used directly only for `pq.NullTime`)
- `golang.org/x/crypto` v0.52.0 (available for `bcrypt` but unused)

### Notes

- `lib/pq` is imported only for `pq.NullTime`. It can be replaced with `pgtype.Timestamptz` from `pgx/v5` to remove an extra dependency.
- `suikat` is pinned; acceptable for an internal framework.

---

## 10. Code Style & Quality

- `go vet` — clean, no issues.
- Naming: Go-idiomatic (`PascalCase` exported, `camelCase` unexported).
- Comments: package-level and public-type docs are present.
- Error wrapping: `ge.Pin()` from suikat + `fmt.Errorf` with `%w`. Sentinel errors (`ErrUserNotFound`, etc.) are used correctly.
- Repetitive code in `authorized_service.go` — every method copies ~10 lines of permission-check logic. Extracting a helper would improve maintainability.

---

## 11. Potential Bugs

| # | Bug | Location | Impact |
|---|-----|----------|--------|
| 1 | **Email uniqueness + soft-delete** | `impl.go:44-49` | `Create` does not check `deleted_at` when looking for duplicate emails. Restoring a soft-deleted user may cause a unique constraint violation. |
| 2 | **`strings.TrimLeft` misused as prefix trim** | `middlewares.auth.go:23` | `TrimLeft(tokenHeader, "Bearer")` removes **any** of the characters `B`, `e`, `a`, `r` from the left. A header like `"Bearererer XXX"` becomes `" XXX"`. Should be `strings.TrimPrefix`. |
| 3 | **Type assertion fragility in logout** | `handlers.auth.go:76` | Assumes `c.Get(Token)` returns `*goxus.UsersToken`. A middleware change could cause a runtime panic. |
| 4 | **`pgx.ErrNoRows` wrapping uncertainty** | `impl.go:206` | Checks `errors.Is(err, pgx.ErrNoRows)`, but if `GetUserByEmail` wraps the error via `ge.Pin`, `errors.Is` may fail. Verify in xo-generated code. |

---

## 12. Recommendations (Priority)

### P0 — Critical (Blockers for Production)

1. - [ ] Hash passwords with `bcrypt`/`argon2`; remove MD5 usage.
2. - [ ] Fix CORS: use a domain whitelist instead of `*`.
3. - [ ] Add HTTPS/TLS termination.
4. - [ ] Replace `strings.TrimLeft` with `strings.TrimPrefix` in auth middleware.

### P1 — High

5. - [x] Add rate limiting on login.
6. - [x] Add expires to `users_tokens` + a cleanup job.
         создай настройку времени жизни токенов = 7 days. 
         создай cronjob для soft_delete устаревших токенов. 
         время жизни измеряется от last_used_at.

7. - [x] Add unique index: `users(email)`.
8. - [x] Replace `GetAll()` with targeted SQL (`JOIN` + `WHERE`) in RBAC checks.
9. - [ ] Add pagination (`limit`/`offset`) to `ListUsers`.
10. - [ ] Switch to structured JSON logging (e.g., `slog` or `zap`).

### P2 — Medium

11. - [ ] Add tests for handlers and middlewares (`httptest` + Gin).
12. - [ ] Add integration tests for Wire initialization.
13. - [ ] Extract a generic RBAC-check helper in `authorized_service.go`.
14. - [ ] Add request ID middleware.
15. - [ ] Remove `log.Fatal` from graceful shutdown.
16. - [ ] Replace `pq.NullTime` with `pgx`/`pgtype` types.

### P3 — Low / Tech Debt

17. - [ ] Add OpenAPI/Swagger documentation.
18. - [ ] Enhance health check with DB connectivity probe.
19. - [ ] Add Prometheus metrics.
20. - [x] Define an API versioning strategy (v1 structure is ready).

---

## 13. Final Score

| Category | Grade | Comment |
|----------|-------|---------|
| Project Structure | **A** | Clear layer separation, proper use of `internal/`. |
| DI / Wire | **A** | Correct, type-safe. |
| Domain Layer / DDD | **B+** | Good auth decorator, but some code duplication. |
| RBAC Design | **B** | Feature-complete, but N+1 / full-scan queries. |
| API Layer | **B-** | REST conventions OK, but missing pagination & rate limiting. |
| Security | **D+** | Plaintext passwords, CORS `*`, no TLS — critical. |
| Testing | **B+** | Excellent domain coverage, but 0% on API & infrastructure. |
| Graceful Shutdown | **B** | Works, but `log.Fatal` is risky. |
| Migrations / DB Schema | **B** | Structure OK, but gaps in soft-delete uniqueness. |
| Code Style / Quality | **A-** | `go vet` clean, idiomatic Go, minor duplication. |

### Overall Grade: **B**

> Good foundational architecture with clear separation of concerns and solid domain-level testing. However, **critical security issues** (plaintext passwords, open CORS, missing TLS) block any production deployment. Once those are addressed, the codebase provides a solid foundation for a SaaS admin panel.
