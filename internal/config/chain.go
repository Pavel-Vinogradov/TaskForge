package config

import (
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type (
	Services struct {
		App        *AppService
		HttpServer *HttpServer
		DB         *sql.DB
		Redis      *redis.Client
		Logger     *logrus.Logger
		Validator  *validator.Validate
		Cache      *CacheService
		UseCases   *UseCases
		Warnings   []string
	}

	union struct {
		Internal *Services
	}

	configurator struct {
		member serviceConfigurator
		union  *union
	}

	serviceConfigurator interface {
		Health() error
		setupMember(s *Services) error
		appendMemberToServices(*Services)
	}

	ServiceChainMember struct {
		serviceConfigurator
	}

	memberHandler func()
)

func NewConfigurations(cfg AppConfig) (s *Services) {
	unions := newUnion()

	handlers := []memberHandler{
		func() {
			runServiceChain([]configurator{
				{member: newLoggerService(), union: unions},
				{member: newValidatorService(), union: unions},
			}...)
		},
		func() {
			runServiceChain([]configurator{
				{member: newAppService(cfg), union: unions},
			}...)
		},
		func() {
			if cfg.DB != nil {
				runServiceChain([]configurator{
					{member: newDatabaseService(cfg.DB), union: unions},
				}...)
			}
		},
		func() {
			if cfg.Redis != nil {
				runServiceChain([]configurator{
					{member: newRedisService(cfg.Redis), union: unions},
				}...)
			}
		},
		func() {
			runServiceChain([]configurator{
				{member: newCacheService(), union: unions},
				{member: newHttpServerService(), union: unions},
			}...)
		},
	}

	return unions.withServicesChainHandler(handlers...)
}

func newUnion() *union {
	return &union{
		Internal: &Services{},
	}
}

func (u union) withServicesChainHandler(handlers ...memberHandler) *Services {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("error while execute one of handlers in servicesChain: %v", r)
		}
	}()

	for _, handler := range handlers {
		handler()
	}

	return u.Internal
}

func runServiceChain(members ...configurator) {
	for _, m := range members {
		err := m.member.setupMember(m.union.Internal)

		if err != nil {
			var appServiceError AppServiceError
			if errors.As(err, &appServiceError) {
				panic(appServiceError)
			}

			m.union.Internal.Warnings = append(m.union.Internal.Warnings, err.Error())
			logrus.Warn(err)
		}

		m.member.appendMemberToServices(m.union.Internal)
	}
}
