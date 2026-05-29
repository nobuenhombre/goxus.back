#======================================================
# goxus
#======================================================
APP_NAME=goxus
APP_LOG_NAME=goxus

PROJECT_ROOT_PATH=../../../..
INSTALL_PATH=/usr/local/bin

BUILD_PLATFORM=linux
BIN_PATH=bin/$(APP_NAME)/linux

APP_BINARY=$(BIN_PATH)/$(APP_NAME)
APP_INSTALL=$(INSTALL_PATH)/$(APP_NAME)

SERVICE_NAME="api_goxus"
SERVICE_PATH=/etc/systemd/system

#=========================================================================
#
#=========================================================================
.PHONY: help build

help: Makefile
	@echo "Выберите опцию сборки "$(BINARY_NAME)":"
	@sed -n 's/^##//p' $< | column -s ':' |  sed -e 's/^/ /'

## all: Скомпилировать приложение и установить
all: uninstall-app build-app install-app

## build-app: Скомпилировать приложение
build-app:
	cd $(PROJECT_ROOT_PATH)/ && \
	go mod tidy && \
	CGO_ENABLED=0 GOOS=$(BUILD_PLATFORM) GOARCH=amd64 go build -ldflags="-s -w" -o $(APP_BINARY) -v ./src/cmd/$(APP_NAME) && \
	chmod +x $(APP_BINARY) && \
	ls -lh $(APP_BINARY) && \
	$(APP_BINARY) --version;

## build-app-progress: Скомпилировать с прогресс-баром
build-app-progress:
	cd $(PROJECT_ROOT_PATH)/ && \
	CGO_ENABLED=0 GOOS=$(BUILD_PLATFORM) GOARCH=amd64 \
	GOFLAGS="-ldflags=-s -w" \
	gobp --full-rebuild -binary ./src/cmd/$(APP_NAME) -out $(APP_BINARY) && \
	chmod +x $(APP_BINARY) && \
	ls -lh $(APP_BINARY) && \
	$(APP_BINARY) --version;

## install-app: Установить приложение
install-app:
	sudo mkdir -p /var/log/$(APP_LOG_NAME) && \
	sudo chmod 777 /var/log/$(APP_LOG_NAME) && \
	sudo ln -sf $(shell pwd)/$(PROJECT_ROOT_PATH)/$(APP_BINARY) $(APP_INSTALL);

## uninstall-app: Удалить приложение
uninstall-app:
	sudo rm -f $(APP_INSTALL);

## install-service: Установить HTTP Сервис
install-service:
	sudo systemctl enable $(shell pwd)/$(PROJECT_ROOT_PATH)/configs/$(SERVER_ROLE)/$(SERVICE_NAME).service && \
	sudo systemctl daemon-reload;

## uninstall-service: Удалить HTTP Сервис
uninstall-service:
	sudo systemctl disable $(shell pwd)/$(PROJECT_ROOT_PATH)/configs/$(SERVER_ROLE)/$(SERVICE_NAME).service && \
	sudo systemctl daemon-reload;

## service-stop: остановить http
service-stop:
	sudo systemctl stop $(SERVICE_NAME);

## service-start: запустить http
service-start:
	sudo systemctl start $(SERVICE_NAME);

## service-status: получить статус http
service-status:
	sudo systemctl status $(SERVICE_NAME);

## service-restart: перезагрузить http
service-restart:
	sudo systemctl restart $(SERVICE_NAME);

## upgrade: upgrade
upgrade: service-stop build-app-progress service-start service-status
