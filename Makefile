PROJECT_NAME=goxus
COVER_THRESHOLD ?= 85.0

export GOPROXY := https://proxy.golang.org,direct
export GOSUMDB := off

#===========================================
# команды makefile -
# если команда совпадет с названием каталога
#===========================================
.PHONY: help deps wire test test-cover test-coverage-threshold

help: Makefile
	@echo "Выберите опцию сборки:"
	@sed -n 's/^##//p' $< | column -s ':' |  sed -e 's/^/ /'

## deps: Инициализация модулей, скачать все необходимые программе модули
deps:
	rm -f go.mod
	rm -f go.sum
	go mod init $(PROJECT_NAME)
	go get -u ./...

## wire: Генерация wire_gen.go через Google Wire
wire:
	wire ./src/...
	gofmt -w src/cmd/*/wire_gen.go

## test: Запуск тестов с race detection
test:
	go test ./... -race -count=1

## test-cover: Запуск тестов с coverage отчётом
test-cover:
	go test ./... -race -count=1 -coverprofile=c.out
	go tool cover -func=c.out
	@rm -f c.out

## test-coverage-threshold: Проверка coverage >= порога (по умолч. $(COVER_THRESHOLD)%)
test-coverage-threshold:
	go test ./... -race -count=1 -coverprofile=c.out
	@COVER=$$(go tool cover -func=c.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "---"; \
	echo "Total coverage: $$COVER%"; \
	if [ "$$(echo "$$COVER < $(COVER_THRESHOLD)" | bc -l)" -eq 1 ]; then \
		echo "❌ FAIL: coverage $$COVER% < $(COVER_THRESHOLD)%"; \
		rm -f c.out; \
		exit 1; \
	else \
		echo "✅ PASS: coverage $$COVER% >= $(COVER_THRESHOLD)%"; \
	fi
	@rm -f c.out
