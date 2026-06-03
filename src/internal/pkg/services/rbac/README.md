# RBAC Service — Role-Based Access Control

[![Go Version](https://img.shields.io/badge/Go-1.26.1-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue)](LICENSE)

**rbac** is a Role-Based Access Control service for the [goxus](https://github.com/nobuenhombre/goxus) backend. It manages roles, permissions, and user-role assignments using PostgreSQL through xo-generated repositories.

---

## Data Model

```
User ──< rbac_user_roles >── Role >── rbac_role_permissions >── Permission
```

| Table | Description | Key Fields |
|-------|-------------|------------|
| `rbac_roles` | Named roles (e.g. "admin", "moderator") | `id`, `name`, `slug` (unique) |
| `rbac_permissions` | Granular actions (e.g. "create_post") | `id`, `name`, `slug` (unique) |
| `rbac_role_permissions` | Role ↔ Permission (M:N) | `role_id`, `permission_id` |
| `rbac_user_roles` | User ↔ Role (M:N) | `user_id`, `role_id` |

Every table includes `id`, `created_at`, `updated_at`.

---

## Quick Start

### Create a role and permission

```go
svc := rbac.New(dbRepo)

// Create
svc.CreateRole("Administrator", "admin")
svc.CreatePermission("Create Post", "create_post")

// Assign
svc.AssignPermissionsToRole("admin", []string{"create_post"})
svc.AssignRoleToUser(42, "admin")

// Check
hasRole, _ := svc.CheckUserRole(42, "admin")             // true
canDo, _ := svc.CheckUserPermission(42, "create_post")   // true

// Revoke
svc.RevokeUserRole(42, "admin")
svc.RevokeRolePermission("admin", "create_post")
```

---

## Service Interface

### Role CRUD

| Method | Description |
|--------|-------------|
| `CreateRole(name, slug string) error` | Create a new role |
| `GetAllRoles() ([]*goxus.RbacRole, error)` | List all roles |
| `GetUserRoles(userID int64) ([]*goxus.RbacRole, error)` | List a user's roles |
| `DeleteRole(slug string) error` | Delete a role (fails if assigned to any user) |

### Permission CRUD

| Method | Description |
|--------|-------------|
| `CreatePermission(name, slug string) error` | Create a new permission |
| `GetAllPermissions() ([]*goxus.RbacPermission, error)` | List all permissions |
| `GetRolePermissions(slug string) ([]*goxus.RbacPermission, error)` | List a role's permissions |
| `DeletePermission(slug string) error` | Delete a permission (fails if assigned to any role) |

### Assignment

| Method | Description |
|--------|-------------|
| `AssignPermissionsToRole(roleSlug string, permSlugs []string) error` | Grant permissions to a role |
| `AssignRoleToUser(userID int64, roleSlug string) error` | Assign a role to a user |

### Revocation

| Method | Description |
|--------|-------------|
| `RevokeUserRole(userID int64, roleSlug string) error` | Remove a role from a user |
| `RevokeRolePermission(roleSlug, permSlug string) error` | Remove a permission from a role |

### Verification

| Method | Description |
|--------|-------------|
| `CheckUserRole(userID int64, slug string) (bool, error)` | Does the user have this role? |
| `CheckUserPermission(userID int64, slug string) (bool, error)` | Does the user have this permission? (resolved via roles) |
| `CheckRolePermission(roleSlug, permSlug string) (bool, error)` | Does the role have this permission? |

---

## Wire DI Integration

```go
// The package exports a Wire provider set:
var rbac.ProviderSet = wire.NewSet(rbac.ProvideRbac)

// Usage in a Wire injector:
wire.Build(
    db.ProviderSet,
    rbac.ProviderSet,
    // ...
)
```

`ProvideRbac` returns `(Service, func(), error)` — the cleanup function logs the service shutdown.

---

## Error Handling

| Error Constant | Meaning |
|----------------|---------|
| `ErrAlreadyExists` | Role or permission slug already exists |
| `ErrRoleNotFound` / `ErrPermissionNotFound` | Slug not found in the database |
| `ErrRoleInUse` / `ErrPermissionInUse` | Cannot delete — still assigned to users/roles |

All errors are wrapped with context (the slug that caused the failure).

---

## Rules & Conventions

- **Slugs** are the primary identifier for roles and permissions — they should be unique and human-readable (e.g. `"admin"`, `"create_post"`)
- **User IDs** are `int64` and reference the application's user table (outside the RBAC package scope)
- **Deletion guard** — roles with active user assignments and permissions with active role assignments cannot be deleted
- **Duplicate-safe** assignment — assigning an already existing role/permission link is silently skipped
- **No RBAC logic outside this package** — all access controls go through the `Service` interface

---

## Schema (PostgreSQL)

```sql
CREATE TABLE rbac_roles (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    slug       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rbac_permissions (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    slug       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rbac_role_permissions (
    id            BIGSERIAL PRIMARY KEY,
    role_id       BIGINT NOT NULL REFERENCES rbac_roles(id),
    permission_id BIGINT NOT NULL REFERENCES rbac_permissions(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rbac_user_roles (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL,
    role_id    BIGINT NOT NULL REFERENCES rbac_roles(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Tech Stack

| Component | Choice |
|-----------|--------|
| Language | Go 1.26 |
| DB Driver | pgx (via `suikat/pkg/db/connectors/postgres-pgx-db`) |
| ORM / Codegen | xo (type-safe Go from SQL) |
| DI | Google Wire |

---

## License

Apache 2.0 — see [LICENSE](../../../../../../LICENSE)
