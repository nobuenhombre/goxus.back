# rbac — Role-Based Access Control Service

## Назначение

Пакет `rbac` предоставляет сервис управления ролями и правами доступа (RBAC). Реализует классическую модель: пользователи → роли → разрешения, где каждому пользователю можно назначить несколько ролей, а каждой роли — несколько разрешений.

- использует pgx DB запросы
- Google Wire DI через `provider.go`
- все возвращаемые ошибки обёрнуты через `ge.Pin()` (пакет `github.com/nobuenhombre/suikat/pkg/ge`)

## Файлы

| Файл | Назначение |
|------|------------|
| `service.go` | Интерфейс `Service` + конкретная реализация `impl` |
| `provider.go` | Google Wire-провайдер `ProvideRbac` |
| `setup_test.go` | TestMain (контейнер + миграции + test user), truncateRBAC, setupTest |
| `service_test.go` | Все тесты — табличные, testify, testcontainers |

## Интерфейс Service

### CRUD ролей и разрешений

| Метод | Описание |
|-------|----------|
| `CreateRole(name, slug)` | Создать новую роль |
| `CreatePermission(name, slug)` | Создать новое разрешение |
| `DeleteRole(slug)` | Удалить роль + очистить role_permissions (только если не назначена пользователям) |
| `DeletePermission(slug)` | Удалить разрешение (только если не назначено ролям) |
| `GetAllRoles()` | Получить все роли |
| `GetAllPermissions()` | Получить все разрешения |
| `GetUserRoles(userID)` | Получить все роли пользователя (фильтрация in-memory из GetAll) |
| `GetRolePermissions(slug)` | Получить все разрешения роли (фильтрация in-memory из GetAll) |

### Назначение и проверка

| Метод | Описание |
|-------|----------|
| `AssignPermissionsToRole(roleSlug, permSlugs)` | Назначить разрешения роли — **идемпотентно**: уже существующие связи пропускаются (continue) |
| `AssignRoleToUser(userID, roleSlug)` | Назначить роль пользователю — если уже назначена, возвращает `ErrAlreadyExists` |
| `CheckUserRole(userID, roleSlug)` | Проверить, есть ли у пользователя роль |
| `CheckUserPermission(userID, permSlug)` | Проверить, есть ли у пользователя разрешение (через все его роли) |
| `CheckRolePermission(roleSlug, permSlug)` | Проверить, есть ли у роли разрешение |

### Отзыв

| Метод | Описание |
|-------|----------|
| `RevokeUserRole(userID, roleSlug)` | Отозвать роль у пользователя — **no-op** если не была назначена |
| `RevokeRolePermission(roleSlug, permSlug)` | Отозвать разрешение у роли — **no-op** если не было назначено |

## Важная семантика дублей

- `AssignPermissionsToRole` — **безопасный пропуск** (continue при существующей связи). Повторный вызов не дублирует строки в `rbac_role_permissions`.
- `AssignRoleToUser` — **возвращает `ErrAlreadyExists`**. Повторное назначение той же роли тому же пользователю — ошибка.
- `AssignRoleToUser` — **всегда требует нового user'а** для каждой новой роли. Если у пользователя уже есть роль `moderator`, назначить ему снова `moderator` нельзя.

## Модель данных (4 таблицы)

```
rbac_roles
  ├─ id: int64, slug: string (unique)
  │
  ├── rbac_role_permissions  (many-to-many: role ↔ permission)
  │     ├─ role_id       → rbac_roles.id
  │     └─ permission_id → rbac_permissions.id
  │
  └── rbac_user_roles  (many-to-many: user ↔ role)
        ├─ user_id (ссылка на users.id вне пакета rbac)
        └─ role_id → rbac_roles.id

rbac_permissions
  └─ id: int64, slug: string (unique)
```

Каждая сущность: `id`, `name`, `slug`, `created_at`, `updated_at`.
Связующие таблицы: `id`, `role_id`/`permission_id`/`user_id`, `created_at`, `updated_at`.

## Wire-интеграция

```go
// provider.go
var ProviderSet = wire.NewSet(ProvideRbac)

func ProvideRbac(dbRepo *goxus.DbGoxusRepo) (Service, func(), error)
```

Провайдер возвращает `(Service, func(), error)` — стандартный контракт Wire для goxus. Cleanup логирует завершение работы сервиса через `log.Println`. Cleanup-функция должна оставаться ненулевой (контракт Wire).

## Ошибки

Все ошибки определены как `var`-переменные, оборачиваются через `ge.Pin(fmt.Errorf(...))`:

| Ошибка | Условие |
|--------|---------|
| `ErrPermissionInUse` | `"cannot delete assigned permission"` |
| `ErrPermissionNotFound` | `"permission not found"` |
| `ErrRoleInUse` | `"cannot delete assigned role"` |
| `ErrRoleNotFound` | `"role not found"` |
| `ErrAlreadyExists` | `"already exists"` |

Ошибки NotFound возвращаются как `fmt.Errorf(...)` без `ge.Pin` в Check-методах (используются для логики `bool, error`). Ошибки в методах модификации всегда обёрнуты через `ge.Pin`.

## Тестирование

### TestMain (setup_test.go)

- Запускает один PostgreSQL-контейнер на все тесты через `testpostgres.StartPostgresContainer`
- Применяет миграции из `../../../../scripts/xo/goxus/migrations/`
- Создаёт динамического тестового пользователя (`goxus.User`) — больше не зависит от seeded ID=1
- Закрывает репозиторий и терминирует контейнер после всех тестов

### setupTest (per-test)

- Вызывает `truncateRBAC()` — очищает все 4 RBAC-таблицы в порядке FK: rbac_role_permissions → rbac_user_roles → rbac_permissions → rbac_roles
- Возвращает `testFixtures{svc, repo}`
- `t.Cleanup` повторно вызывает `truncateRBAC()` для изоляции
- `noopLog` — заглушка для `WriteLog` (паникует на nil Log)

### service_test.go

- testify `assert` + `require`
- Каждый тест изолирован через `setupTest(t)`
- Покрытие: создание, дубликаты, назначение, проверка (все варианты), отзыв (включая no-op), удаление (все сценарии), пустые списки, идемпотентность AssignPermissionsToRole, удаление роли с разрешениями

### Запуск

```bash
cd back/src
go test ./internal/pkg/services/rbac/ -v -count=1
```

## Конвенции

- Конструктор `New(dbRepo *goxus.DbGoxusRepo) Service` — экспортирован, принимает xo-репозиторий
- `noopLog` — заглушка для `WriteLog` (паникует на nil, поэтому всегда передаём noopLog)
- При добавлении новой таблицы RBAC в БД — добавить xo-репозиторий в `DbGoxusRepo` и соответствующий метод в `Service` / `impl`
- При изменении логики проверки прав — обновить оба метода `CheckUserRole` и `CheckUserPermission` для консистентности
- Не хранить RBAC-логику вне этого пакета — обращаться через интерфейс `Service`
- `go-sumtype` не используется — RBAC типы не входят в go-sumtype exhaustive checks
