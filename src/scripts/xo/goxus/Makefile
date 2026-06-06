LINT=../lint.sh
BACKUP=../backup.sh
RESTORE=../restore.sh
CREATE=../create.sh
RESET=../reset.sh
MIGRATE_NEW=../migrate-new.sh
MIGRATE_UP=../migrate-up.sh
MIGRATE_DOWN=../migrate-down.sh

#===========================================
# команды makefile -
# если команда совпадет с названием каталога
#===========================================
.PHONY: help gen backup backup-production restore lint create create-e2e reset reset-e2e migrate-new migrate-up migrate-up-e2e migrate-down migrate-down-e2e

help: Makefile
	@echo "Выберите опцию сборки:"
	@sed -n 's/^##//p' $< | column -s ':' |  sed -e 's/^/ /'

## gen: сгенерировать код GO
gen:
	goxogen -config "xo.yaml"

## backup: сделать бекап локальной базы
backup:
	time $(BACKUP) xo.yaml

## backup-production: сделать бекап продакшн базы
backup-production:
	time $(BACKUP) xo.production.yaml

## restore: восстановить бекап в локальную базы
restore:
	time $(RESTORE) xo.yaml

## lint: проверить код
lint:
	time $(LINT) xo.yaml

## create: создать локальную базу
create:
	time $(CREATE) xo.yaml

## create-e2e: создать локальную базу e2e
create-e2e:
	time $(CREATE) xo.e2e.yaml

## reset: reset локальную базу
reset:
	time $(RESET) xo.yaml

## reset-e2e: reset локальную базу e2e
reset-e2e:
	time $(RESET) xo.e2e.yaml

## migrate-new: создать новую миграцию БД
migrate-new:
	$(MIGRATE_NEW) migrations

## migrate-up: накатить миграции в локальную базу
migrate-up:
	time $(MIGRATE_UP) xo.yaml

## migrate-up-e2e: накатить миграции в e2e базу
migrate-up-e2e:
	time $(MIGRATE_UP) xo.e2e.yaml

## migrate-down: откатить последнюю миграцию в локальной базе
migrate-down:
	time $(MIGRATE_DOWN) xo.yaml

## migrate-down-e2e: откатить последнюю миграцию в e2e базе
migrate-down-e2e:
	time $(MIGRATE_DOWN) xo.e2e.yaml
