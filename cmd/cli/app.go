package cli

import (
	"TaskForge/internal/config"
	"TaskForge/internal/handler"
	"TaskForge/internal/infrastructure/repository"
	"TaskForge/internal/middleware"
	"TaskForge/internal/usecase"
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	UseCases    *config.UseCases
	Middlewares *config.Middlewares
	Services    *config.Services
}

func NewApp(servicesConfig *config.Services) (*App, error) {
	authRepo := repository.NewAuthRepository(servicesConfig.DB)

	useCases := &config.UseCases{
		Auth: usecase.NewAuthUseCase(authRepo),
	}

	mws := &config.Middlewares{
		Cors: middleware.CorsMiddleware{},
		JWT:  middleware.JWTAuthMiddleware{},
	}

	app := &App{
		UseCases:    useCases,
		Services:    servicesConfig,
		Middlewares: mws,
	}

	return app, nil
}

func (app *App) RunApi(ctx context.Context) error {
	r := gin.Default()

	r.Use(app.Middlewares.Cors.CorsMiddleware)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authHandler := handler.NewAuthHandler(app.UseCases.Auth)
	api := r.Group("/api/v1")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(app.Middlewares.JWT.JWTAuthMiddleware())
		{
			//api.POST("/teams", app.UseCases.Teams.CreateHandler)
			//api.GET("/teams", app.UseCases.Teams.ListHandler)
			//api.POST("/teams/:id/invite", app.UseCases.Teams.InviteHandler)
			//
			//api.POST("/tasks", app.UseCases.Tasks.CreateHandler)
			//api.GET("/tasks", app.UseCases.Tasks.ListHandler)
			//api.PUT("/tasks/:id", app.UseCases.Tasks.UpdateHandler)
			//api.GET("/tasks/:id/history", app.UseCases.Tasks.HistoryHandler)
		}

	}

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(app.Services.HttpServer.Port),
		Handler: r,
	}

	errCh := make(chan error, 1)

	go func() {
		logrus.Infof("Starting REST API server on port %d", app.Services.HttpServer.Port)
		logrus.Infof("Swagger documentation available at: http://localhost:%d/swagger/index.html", app.Services.HttpServer.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}
