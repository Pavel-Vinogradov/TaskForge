package config

import (
	"TaskForge/internal/interfaces/auth"
	"TaskForge/internal/interfaces/task"
	"TaskForge/internal/interfaces/team"

	"TaskForge/internal/middleware"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

// AppConfig - конфигурация приложения
type AppConfig struct {
	DB          *DBConfig
	Redis       *RedisConfig
	JWT         *JWTConfig
	Middlewares *Middlewares
}

// AppService - контейнер для конфигурации и use cases
type AppService struct {
	Config   AppConfig
	UseCases *UseCases
}

// UseCases - бизнес-логика
type UseCases struct {
	Auth  auth.UseCaseAuth
	Teams team.UseCaseTeam
	Tasks task.UseCaseTask
}

type Middlewares struct {
	Cors middleware.CorsMiddleware
	JWT  *middleware.JWTAuthMiddleware
}

type AppServiceError string

func (e AppServiceError) Error() string {
	return fmt.Sprintf("app service error: %s", string(e))
}

// DBConfig - настройки базы данных
type DBConfig struct {
	Driver       string
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLife  time.Duration
}

// RedisConfig - настройки Redis
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	TTL      time.Duration
}

// JWTConfig - настройки JWT
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// HttpServer - HTTP сервер
type HttpServer struct {
	Engine *gin.Engine
	Port   int
}

// CacheService - сервис кеширования
type CacheService struct {
	Redis *redis.Client
	TTL   time.Duration
}

func NewAppService(cfg AppConfig, useCases *UseCases) *AppService {
	return &AppService{
		Config:   cfg,
		UseCases: useCases,
	}
}

// Health проверка состояния сервисов
func (a *AppService) Health(ctx context.Context, dbPing func(c context.Context) error, redisPing func(c context.Context) error) error {

	var errs []error

	if err := dbPing(ctx); err != nil {
		errs = append(errs, fmt.Errorf("db ping failed: %w", err))
	}
	if err := redisPing(ctx); err != nil {
		errs = append(errs, fmt.Errorf("redis ping failed: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("health check errors: %v", errs)
	}

	return nil
}
