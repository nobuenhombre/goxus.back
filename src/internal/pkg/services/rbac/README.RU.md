# RBAC Service — Контроль доступа на основе ролей

[![Go Version](https://img.shields.io/badge/Go-1.26.1-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue)](LICENSE)

**rbac** — сервис контроля доступа на основе ролей (Role-Based Access Control) для бэкенда [goxus](https://github.com/nobuenhombre/goxus). Управляет ролями, разрешениями и назначением ролей пользователям через PostgreSQL с использованием xo-сгенерированных репозиториев.

---

## Модель данных

```
Пользователь ──< rbac_user_roles >── Роль >── rbac_role_permissions >── Разрешение
```

| Таблица | Описание | Ключевые поля |
|---------|----------|---------------|
| `rbac_roles` | Поименованные роли (например "admin", "moderator") | `id`, `name`, `slug` (уникальный) |
| `rbac_permissions` | Конкретные действия ("create_post") | `id`, `name`, `slug` (уникальный) |
| `rbac_role_permissions` | Связь Роль ↔ Разрешение (M:N) | `role_id`, `permission_id` |
| `rbac_user_roles` | Связь Пользователь ↔ Роль (M:N) | `user_id`, `role_id` |

Все таблицы содержат `id`, `created_at`, `updated_at`.

---

## Быстрый старт

### Создание роли и разрешения

```go
svc := rbac.New(dbRepo)

// Создание
svc.CreateRole("Administrator", "admin")
svc.CreatePermission("Create Post", "create_post")

// Назначение
svc.AssignPermissionsToRole("admin", []string{"create_post"})
svc.AssignRoleToUser(42, "admin")

// Проверка
hasRole, _ := svc.CheckUserRole(42, "admin")             // true
canDo, _ := svc.CheckUserPermission(42, "create_post")   // true

// Отзыв
svc.RevokeUserRole(42, "admin")
svc.RevokeRolePermission("admin", "create_post")
```

---

## Интерфейс Service

### Управление ролями

| Метод | Описание |
|-------|----------|
| `CreateRole(name, slug string) error` | Создать новую роль |
| `GetAllRoles() ([]*goxus.RbacRole, error)` | Получить все роли |
| `GetUserRoles(userID int64) ([]*goxus.RbacRole, error)` | Получить роли пользователя |
| `DeleteRole(slug string) error` | Удалить роль (только если не назначена) |

### Управление разрешениями

| Метод | Описание |
|-------|----------|
| `CreatePermission(name, slug string) error` | Создать новое разрешение |
| `GetAllPermissions() ([]*goxus.RbacPermission, error)` | Получить все разрешения |
| `GetRolePermissions(slug string) ([]*goxus.RbacPermission, error)` | Получить разрешения роли |
| `DeletePermission(slug string) error` | Удалить разрешение (только если не назначено) |

### Назначение

| Метод | Описание |
|-------|----------|
| `AssignPermissionsToRole(roleSlug string, permSlugs []string) error` | Выдать разрешения роли |
| `AssignRoleToUser(userID int64, roleSlug string) error` | Назначить роль пользователю |

### Отзыв

| Метод | Описание |
|-------|----------|
| `RevokeUserRole(userID int64, roleSlug string) error` | Отозвать роль у пользователя |
| `RevokeRolePermission(roleSlug, permSlug string) error` | Отозвать разрешение у роли |

### Проверка прав

| Метод | Описание |
|-------|----------|
| `CheckUserRole(userID int64, slug string) (bool, error)` | Есть ли у пользователя роль? |
| `CheckUserPermission(userID int64, slug string) (bool, error)` | Есть ли у пользователя разрешение? (через роли) |
| `CheckRolePermission(roleSlug, permSlug string) (bool, error)` | Есть ли у роли разрешение? |

---

## Wire DI Интеграция

```go
// Пакет экспортирует Wire provider set:
var rbac.ProviderSet = wire.NewSet(rbac.ProvideRbac)

// Использование в Wire injector:
wire.Build(
    db.ProviderSet,
    rbac.ProviderSet,
    // ...
)
```

`ProvideRbac` возвращает `(Service, func(), error)` — функция cleanup логирует завершение сервиса.

---

## Обработка ошибок

| Константа ошибки | Значение |
|------------------|----------|
| `ErrAlreadyExists` | Роль или разрешение с таким slug уже существует |
| `ErrRoleNotFound` / `ErrPermissionNotFound` | Slug не найден в БД |
| `ErrRoleInUse` / `ErrPermissionInUse` | Нельзя удалить — ещё назначено пользователям/ролям |

Все ошибки оборачиваются с указанием slug, вызвавшего ошибку.

---

## Правила и конвенции

- **Slug** — первичный идентификатор ролей и разрешений. Должен быть уникальным и человекочитаемым (например `"admin"`, `"create_post"`)
- **ID пользователей** — `int64`, ссылаются на таблицу пользователей приложения (вне скона пакета rbac)
- **Защита от удаления** — роль с активными назначениями пользователям и разрешение с назначением ролям не могут быть удалены
- **Безопасное дублирование** — повторное назначение существующей связи роль-разрешение игнорируется без ошибки
- **Вся RBAC-логика в одном пакете** — контроль доступа только через интерфейс `Service`

---

## Схема (PostgreSQL)

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

| Компонент | Выбор |
|-----------|-------|
| Язык | Go 1.26 |
| Драйвер БД | pgx (через `suikat/pkg/db/connectors/postgres-pgx-db`) |
| ORM / Генерация | xo (типобезопасный Go из SQL) |
| DI | Google Wire |

---

## Лицензия

Apache 2.0 — см. [LICENSE](../../../../../../LICENSE)
