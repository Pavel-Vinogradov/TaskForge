include .env
export DB_USER DB_PASSWORD DB_HOST DB_PORT DB_NAME

MYSQL_URL := mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

.PHONY: help build run test clean migrate-up migrate-down swagger docker-build docker-run docker-up docker-down docker-rebuild docker-logs

# Цель по умолчанию
help:
	@echo "Доступные команды:"
	@echo "  build        - Собрать приложение"
	@echo "  run          - Запустить приложение"
	@echo "  test         - Запустить тесты"
	@echo "  clean        - Очистить артефакты сборки"
	@echo "  migrate-up   - Накатить миграции базы данных"
	@echo "  migrate-down - Откатить миграции базы данных"
	@echo "  swagger      - Сгенерировать Swagger документацию"
	@echo "  docker-build - Собрать Docker образ"
	@echo "  docker-run   - Запустить Docker контейнер"
	@echo "  docker-up    - Запустить через docker-compose"
	@echo "  docker-down  - Остановить docker-compose"
	@echo "  docker-rebuild - Пересобрать и запустить docker-compose"
	@echo "  docker-logs  - Посмотреть логи приложения"

# Собрать приложение
build:
	go build -o bin/taskforge cmd/server/main.go

# Запустить приложение
run:
	go run cmd/server/main.go

# Запустить тесты
test:
	go test -v ./...

# Очистить артефакты сборки
clean:
	rm -rf bin/

# Накатить миграции базы данных
migrate-up:
	migrate -path migrations -database "$(MYSQL_URL)" up

# Откатить миграции базы данных
migrate-down:
	migrate -path migrations -database "$(MYSQL_URL)" down

# Сгенерировать Swagger документацию
swagger:
	swag init -g cmd/server/main.go -o docs

# Собрать Docker образ
docker-build:
	docker build -t taskforge .

# Запустить Docker контейнер
docker-run:
	docker run -p 8080:8080 --env-file .env taskforge

# Запустить через docker-compose
docker-up:
	docker-compose up -d

# Остановить docker-compose
docker-down:
	docker-compose down

# Пересобрать и запустить через docker-compose
docker-rebuild:
	docker-compose up -d --build

# Посмотреть логи docker-compose
docker-logs:
	docker-compose logs -f app

# Процесс разработки
dev: swagger run

# Установить зависимости
deps:
	go mod download
	go mod tidy

# Проверить код линтером
lint:
	golangci-lint run

# Форматировать код
fmt:
	go fmt ./...
