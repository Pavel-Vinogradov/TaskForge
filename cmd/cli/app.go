package cli

import (
	"TaskForge/internal/config"
	"TaskForge/internal/handler"
	"TaskForge/internal/infrastructure/repository"
	"TaskForge/internal/middleware"
	"TaskForge/internal/usecase"
	"TaskForge/pkg/jwt"
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
	teamRepo := repository.NewTeamRepository(servicesConfig.DB)
	taskRepo := repository.NewTaskRepository(servicesConfig.DB)
	taskHistoryRepo := repository.NewTaskHistoryRepository(servicesConfig.DB)

	jwtManager := jwt.NewManager(
		servicesConfig.App.Config.JWT.Secret,
		servicesConfig.App.Config.JWT.Expiration,
	)

	taskUC := usecase.NewTaskUseCase(taskRepo, taskHistoryRepo)
	historyObserver := usecase.NewTaskHistoryObserver(taskHistoryRepo)
	taskUC.AddObserver(historyObserver)

	useCases := &config.UseCases{
		Auth:  usecase.NewAuthUseCase(authRepo),
		Teams: usecase.NewTeamUseCase(teamRepo),
		Tasks: taskUC,
	}

	mws := &config.Middlewares{
		Cors:        middleware.CorsMiddleware{},
		JWT:         middleware.NewJWTAuthMiddleware(jwtManager),
		UserContext: middleware.UserContextMiddleware{},
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

	jwtManager := jwt.NewManager(
		app.Services.App.Config.JWT.Secret,
		app.Services.App.Config.JWT.Expiration,
	)

	authHandler := handler.NewAuthHandler(app.UseCases.Auth, jwtManager)
	taskHandler := handler.NewTaskHandler(app.UseCases.Tasks)
	teamHandler := handler.NewTeamHandler(app.UseCases.Teams)
	api := r.Group("/api/v1")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(app.Middlewares.JWT.JWTAuthMiddleware())
		protected.Use(app.Middlewares.UserContext.UserContextMiddleware())
		{
			protected.POST("/teams", teamHandler.CreateTeam)
			protected.GET("/teams", teamHandler.ListTeams)
			protected.POST("/teams/:id/invite", teamHandler.InviteUser)

			protected.POST("/tasks", taskHandler.CreateTask)
			protected.GET("/tasks", taskHandler.ListTask)
			protected.PUT("/tasks/:id", taskHandler.UpdateTask)
			protected.GET("/tasks/:id/history", taskHandler.HistoryTask)
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
