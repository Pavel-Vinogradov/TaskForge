package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// Logger Service
type loggerService struct {
	*logrus.Logger
}

func newLoggerService() *ServiceChainMember {
	return &ServiceChainMember{
		serviceConfigurator: &loggerService{Logger: logrus.New()},
	}
}

func (s *loggerService) Health() error {
	return nil
}

func (s *loggerService) setupMember(_ *Services) error {
	s.Logger.SetLevel(logrus.InfoLevel)
	s.Logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	return nil
}

func (s *loggerService) appendMemberToServices(services *Services) {
	services.Logger = s.Logger
}

// Validator Service
type validatorService struct {
	*validator.Validate
}

func newValidatorService() *ServiceChainMember {
	return &ServiceChainMember{
		serviceConfigurator: &validatorService{Validate: validator.New()},
	}
}

func (s *validatorService) Health() error {
	return nil
}

func (s *validatorService) setupMember(_ *Services) error {
	return nil
}

func (s *validatorService) appendMemberToServices(services *Services) {
	services.Validator = s.Validate
}

// App Service
type appService struct {
	*AppService
}

func newAppService(cfg AppConfig) *ServiceChainMember {
	return &ServiceChainMember{
		serviceConfigurator: &appService{AppService: NewAppService(cfg, &UseCases{})},
	}
}

func (s *appService) Health() error {
	return nil
}

func (s *appService) setupMember(_ *Services) error {
	return nil
}

func (s *appService) appendMemberToServices(services *Services) {
	services.App = s.AppService
	services.UseCases = s.AppService.UseCases
}

// Database Service
type databaseService struct {
	*sql.DB
	config *DBConfig
}

func newDatabaseService(cfg *DBConfig) *ServiceChainMember {
	return &ServiceChainMember{
		serviceConfigurator: &databaseService{config: cfg},
	}
}

func (s *databaseService) Health() error {
	if s.DB == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.PingContext(ctx)
}

func (s *databaseService) setupMember(_ *Services) error {
	db, err := sql.Open(s.config.Driver, s.config.DSN)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(s.config.MaxOpenConns)
	db.SetMaxIdleConns(s.config.MaxIdleConns)
	db.SetConnMaxLifetime(s.config.ConnMaxLife)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	s.DB = db
	return nil
}

func (s *databaseService) appendMemberToServices(services *Services) {
	services.DB = s.DB
}

// Redis Service
type redisService struct {
	*redis.Client
	config *RedisConfig
}

func newRedisService(cfg *RedisConfig) *ServiceChainMember {
	return &ServiceChainMember{
		serviceConfigurator: &redisService{config: cfg},
	}
}

func (s *redisService) Health() error {
	if s.Client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.Ping(ctx).Err()
}

func (s *redisService) setupMember(_ *Services) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     s.config.Addr,
		Password: s.config.Password,
		DB:       s.config.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return err
	}

	s.Client = rdb
	return nil
}

func (s *redisService) appendMemberToServices(services *Services) {
	services.Redis = s.Client
}

// Cache Service
type cacheService struct {
	*CacheService
}

func newCacheService() *ServiceChainMember {
	return &ServiceChainMember{
		serviceConfigurator: &cacheService{},
	}
}

func (s *cacheService) Health() error {
	return nil
}

func (s *cacheService) setupMember(services *Services) error {
	ttl := 5 * time.Minute
	if services.App.Config.Redis != nil {
		ttl = services.App.Config.Redis.TTL
		if ttl == 0 {
			ttl = 5 * time.Minute
		}
	}

	s.CacheService = &CacheService{
		Redis: services.Redis,
		TTL:   ttl,
	}
	return nil
}

func (s *cacheService) appendMemberToServices(services *Services) {
	services.Cache = s.CacheService
}

// HTTP Server Service
type httpServerService struct {
	*HttpServer
}

func newHttpServerService() *ServiceChainMember {
	return &ServiceChainMember{
		serviceConfigurator: &httpServerService{},
	}
}

func (s *httpServerService) Health() error {
	return nil
}

func (s *httpServerService) setupMember(_ *Services) error {
	port := 8080
	s.HttpServer = &HttpServer{
		Engine: gin.New(),
		Port:   port,
	}

	s.HttpServer.Engine.Use(gin.Logger())
	s.HttpServer.Engine.Use(gin.Recovery())
	return nil
}

func (s *httpServerService) appendMemberToServices(services *Services) {
	services.HttpServer = s.HttpServer
}
