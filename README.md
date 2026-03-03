# TaskForge - Сервис управления задачами

Техническое задание: Сервис управления задачами с командной работой и историей изменений
## Стек технологий

- **Go 1.25.0** - основной язык программирования
- **MySQL 8.0** - основная база данных
- **Redis** - кеширование
- **Docker & Docker Compose** - контейнеризация
- **JWT** - аутентификация
- **Testcontainers** - интеграционные тесты

## Структура проекта

```
TaskForge/
├── cmd/                    # Точки входа приложения
│   ├── server/            # HTTP сервер
│   └── cli/               # CLI утилиты
├── internal/              # Внутренний код приложения
│   ├── config/            # Конфигурация
│   ├── domain/            # Доменная модель
│   │   ├── entity/        # Сущности
│   │   └── repos/         # Интерфейсы репозиториев
│   ├── handler/           # HTTP обработчики
│   ├── infrastructure/    # Инфраструктурный слой
│   │   └── repository/    # Реализация репозиториев
│   ├── middleware/        # Middleware
│   ├── usecase/           # Бизнес-логика
│   └── interfaces/        # Интерфейсы
├── migrations/            # Миграции базы данных
├── configs/              # Конфигурационные файлы
├── docs/                 # Swagger документация
├── docker-compose.yml    # Docker Compose конфигурация
├── Dockerfile           # Docker образ
├── Makefile             # Сборочные скрипты
└── go.mod               # Go модули
```

## Быстрый старт

### 1. Клонирование и установка зависимостей

```bash
git clone git@github.com:Pavel-Vinogradov/TaskForge.git
cd TaskForge
go mod download
```

### 2. Настройка окружения

Скопируйте файл окружения и настройте его:

```bash
cp .env.example .env
```

Отредактируйте `.env` файл:

```env
# База данных
DB_HOST=localhost
DB_PORT=3306
DB_USER=taskforge
DB_PASSWORD=password
DB_NAME=taskforge

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRES_IN=24h

# Сервер
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
```

### 3. Запуск с Docker Compose (рекомендуется)

```bash
# Запуск всех сервисов (MySQL, Redis, приложение)
docker-compose up -d

# Накат миграций
make migrate-up

# Просмотр логов
docker-compose logs -f app
```

### 4. Запуск без Docker

#### 4.1. Запуск MySQL и Redis

```bash
# Запуск MySQL
docker run -d --name mysql-taskforge \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=taskforge \
  -e MYSQL_USER=taskforge \
  -e MYSQL_PASSWORD=password \
  -p 3306:3306 \
  mysql:8.0

# Запуск Redis
docker run -d --name redis-taskforge \
  -p 6379:6379 \
  redis:alpine
```

#### 4.2. Накат миграций

```bash
make migrate-up
```

#### 4.3. Запуск приложения

```bash
# Разработка
make run

# Или сборка и запуск
make build
./bin/taskforge
```

## Доступные команды

```bash
# Сборка приложения
make build

# Запуск приложения
make run

# Запуск всех тестов
make test

# Запуск только интеграционных тестов
make test-integration

# Накат миграций
make migrate-up

# Откат миграций
make migrate-down

# Генерация Swagger документации
make swagger

# Запуск через Docker Compose
make docker-up

# Остановка Docker Compose
make docker-down

# Пересборка и запуск Docker Compose
make docker-rebuild

# Просмотр логов
make docker-logs

# Установка зависимостей
make deps

# Проверка кода линтером
make lint

# Форматирование кода
make fmt
```

## API Документация

После запуска приложения Swagger документация доступна по адресу:
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **JSON**: http://localhost:8080/swagger/doc.json

## Основные эндпоинты
![img.png](img.png)
