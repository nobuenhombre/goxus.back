LINT=../lint.sh
BACKUP=../backup.sh
RESTORE=../restore.sh
CREATE=../create.sh

#===========================================
# команды makefile -
# если команда совпадет с названием каталога
#===========================================
.PHONY: help gen

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
