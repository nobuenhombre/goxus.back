# Goxus

[![Go Version](https://img.shields.io/badge/Go-1.26.1-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue)](LICENSE)

**Goxus** — это бэкенд-приложение на Go с внедрением зависимостей (Google Wire), HTTP API сервером на Gin, планировщиком cron-задач и генерацией кода PostgreSQL через xo. Предназначен для продакшен-развёртывания как systemd-сервис на Linux.

---

## Быстрый старт

```bash
# Сборка бинарника
cd service/deployments/goxus/linux
make build-app

# Запуск в режиме по умолчанию
./bin/goxus/linux/goxus

# Запуск как HTTP-сервис
./bin/goxus/linux/goxus -runtype=service -config=configs/local/config.yaml

# Версия
./bin/goxus/linux/goxus --version

# Перегенерация Wire DI (после изменений в provider'ах)
make wire
```

---

## Возможности

- **3 режима запуска** — `init`, `service` (HTTP + cron), и режим по умолчанию (бизнес-логика)
- **Внедрение зависимостей** через Google Wire — чистый паттерн provider/cleanup
- **HTTP-сервер на Gin** с graceful shutdown (SIGINT/SIGTERM)
- **Планировщик cron-задач** — расписание из YAML конфига
- **YAML-конфигурация** — пресеты для разных окружений (local, develop, production)
- **Генерация кода PostgreSQL** — xo с кастомными Go-шаблонами
- **systemd-интеграция** — установка, запуск, остановка, обновление как сервиса
- **Статический бинарник** — CGO_ENABLED=0, полностью статическая сборка linux/amd64
- **Логирование в файл** — опционально, с перенаправлением стандартного логгера

---

## Архитектура

```
src/
  cmd/goxus/                  # Точка входа
    main.go                   #   Panic recovery, --version fast path, Wire init
    app.go                    #   Оркестратор IApp, переключение режимов
    wire.go                   #   Wire-инжектор (build tag: wireinject)
    wire_gen.go               #   Авто-генерация `make wire`
  internal/
    app/goxus/
      cli/                    # Парсинг CLI-флагов
      config/                 # YAML конфигурация (загрузка/сохранение)
      domain/                 # Сервис бизнес-логики
      log/                    # Управление лог-файлом
      api/server/             # Gin HTTP сервер
        router/               #   Маршруты, хендлеры, middleware
      cron-job/               # Планировщик cron
        jobs/example/         #   Пример cron-задачи
    pkg/
      db/goxus/               # xo-сгенерированные типы PostgreSQL
  scripts/xo/                 # Инструменты для БД (codegen, миграции, бекапы)
configs/                      # Конфиги для разных окружений + systemd unit'ы
  local/                      #   Локальная разработка
  develop/                    #   Dev-сервер
  production/                 #   Продакшен
service/deployments/          # Makefile'ы сборки и деплоя
bin/                          # Скомпилированные бинарники
```

### Цепочка DI

```
main()
  └─ initializeApp() [Wire]
       ├─ ProvideCLI()        → cli.Service
       ├─ ProvideConfigApp()  → configapp.Service      [зависит от CLI]
       ├─ ProvideDomain()     → DomainService           [зависит от CLI + Config]
       ├─ ProvideLogFile()    → ILogFile                [зависит от CLI]
       ├─ ProvideAPI()        → IHTTPServer             [зависит от Config + Log + Domain]
       ├─ ProvideExampleJobs()→ *cron.Cron              [зависит от Domain]
       └─ newApp()            → IApp                    [главный оркестратор]
```

Каждый provider возвращает `(Service, func(), error)` — Wire отслеживает cleanup в обратном порядке создания.

---

## Режимы запуска

| Режим     | Флаг                     | Поведение                                                   |
|-----------|--------------------------|-------------------------------------------------------------|
| `init`    | `-runtype=init`          | Заглушка инициализации (по умолчанию)                       |
| `service` | `-runtype=service`       | Запуск cron + HTTP, graceful shutdown                       |
| default   | (без флага или другой)   | Запуск `domain.Run()` (заглушка бизнес-логики)              |

### CLI флаги

```
-runtype=init|service    Режим запуска (по умолчанию: init)
-config=<путь>           Путь к YAML конфигу (по умолчанию: config.yaml)
-log=<путь>              Путь к лог-файлу (опционально, пусто = stdout)
--version                Вывод версии и выход
```

---

## Конфигурация

### Структура YAML

```yaml
hosts:
  api:
    host: 127.0.0.1      # Адрес привязки
    post: "8080"         # Порт (внимание: "post" а не "port" — см. Гвозди)
cron:
  example_job:
    enabled: true
    schedule: '@every 10m'
```

### Пресеты окружений

| Окружение   | Файл конфига | systemd unit |
|---|---|---|
| Local | `configs/local/config.yaml` | `api_goxus.service` |
| Develop | `configs/develop/config.yaml` | `api_goxus_dev.service` |
| Production | `configs/production/config.yaml` | `api_goxus.service` |

---

## Сборка и деплой

### Разработка (корневой Makefile)

```bash
make wire          # Перегенерация wire_gen.go
make deps          # Переинициализация go.mod и загрузка зависимостей
make help          # Список доступных команд
```

### Продакшен-сборка (service/deployments/goxus/linux/)

```bash
make build-app           # Статический linux/amd64 бинарник
make build-app-progress  # Сборка с прогресс-баром (требуется gawk)
make install-app         # Симлинк в /usr/local/bin/goxus
make uninstall-app       # Удаление симлинка
make install-service     # Активация systemd unit'а (задать SERVER_ROLE)
make uninstall-service   # Деактивация systemd unit'а
make service-start       # systemctl start api_goxus
make service-stop        # systemctl stop api_goxus
make service-restart     # systemctl restart api_goxus
make service-status      # systemctl status api_goxus
make upgrade             # stop → build → start → status
```

---

## База данных

Генерация кода PostgreSQL через [xo](https://github.com/xo/xo) с кастомными Go-шаблонами. Миграции БД — через [golang-migrate/migrate](https://github.com/golang-migrate/migrate) (CLI `migrate`).

```bash
# Перегенерация типов БД из схемы
cd src/scripts/xo/
./xo.sh goxus/xo.yaml

# Миграции (golang-migrate)
cd src/scripts/xo/goxus/
./migrate-new.sh           # Создать новую миграцию
./migrate-up.sh            # Применить миграции
./migrate-down.sh          # Откатить последнюю миграцию
```

---

## Технологический стек

| Компонент | Выбор | Версия |
|---|---|---|
| Язык | Go | 1.26.1 |
| HTTP | Gin | v1.12.0 |
| DI | Google Wire | v0.7.0 |
| CLI флаги | suikat/pkg/clivar | v0.0.170 |
| Cron | robfig/cron/v3 | v3.0.1 |
| Конфиги | gopkg.in/yaml.v3 | v3.0.1 |
| Генерация БД | xo | (кастомные шаблоны) |
| Драйвер БД | pgx | (сгенерирован) |
| Миграции | golang-migrate/migrate | CLI |
| Платформа | Linux | amd64 |
| Сборка | CGO_ENABLED=0 | статическая |

---

## Конвенции

- **Ошибки:** Обёртка через `ge.Pin(err)` из `suikat/pkg/ge`
- **DI провайдеры:** Каждый `ProvideXxx` возвращает `(Service, func(), error)` — cleanup никогда не nil
- **Версия:** Единственный источник истины в `version/version.go` (SemVer)
- **Тестирование:** Table-driven тесты с YAML-фикстурами
- **Git keepfiles:** `.gitkeep` во всех директориях для сохранения структуры в VCS

---

## Гвозди (Gotchas)

- **Ключ port в YAML:** В конфиге используется `post:` (опечатка), не `port:`. Структурный тег в `config-server.go` — `yaml:"post,omitempty"`. Изменение сломает все конфиги окружений.
- **Wire cleanup:** Все provider'ы должны возвращать ненулевой cleanup. Cleanup вызывается в обратном порядке provider'ов.
- **CGO_ENABLED=0:** Сборка принудительно отключает CGO для статической линковки. Если нужны C-библиотеки — править `go-build.mk`.
- **GOSUMDB=off:** База контрольных сумм отключена в Makefile. Не включать без понимания последствий.
- **Корневой vs деплой Makefile:** Корневой `Makefile` минимален (wire + deps). Полная сборка/деплой — в `service/deployments/goxus/linux/Makefile`.
- **xo-сгенерированный код:** `src/internal/pkg/db/goxus/xo_db.xo.go` авто-генерируется — править шаблоны, а не результат.
- **Нет тестов** для пакетов domain, CLI и server — только `config-app_test.go`.
- **Go proxy:** Используется `proxy.golang.org,direct` — для изолированной среды выставить `GOPROXY=off`.

---

## Статус проекта

Текущая версия: **v0.1.0** — начальный каркас с функциональным HTTP-сервером, планировщиком cron и пайплайном генерации БД. Готов к наполнению бизнес-логикой.

---

## Лицензия

Apache 2.0 — см. [LICENSE](LICENSE)