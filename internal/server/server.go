package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/songfei1983/play-go-api/internal/app"
	"github.com/songfei1983/play-go-api/internal/handler"
	mymiddleware "github.com/songfei1983/play-go-api/internal/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type Server struct {
	app    *app.App
	router *echo.Echo
	server *http.Server
}

func New(app *app.App) *Server {
	e := echo.New()
	e.Debug = true
	e.Use(middleware.Logger())
	e.Use(otelecho.Middleware("api-service"))
	e.Use(middleware.RequestID())
	e.Use(mymiddleware.MetricsMiddleware())

	s := &Server{
		app:    app,
		router: e,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", "8080"),
			Handler: e,
		},
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	userHandler := handler.NewUserHandler(s.app.DB, s.app.Redis)

	// API routes
	v1 := s.router.Group("/api/v1")
	v1.POST("/users", userHandler.Register)
	v1.GET("/users/:id", userHandler.GetUser)

	// Metrics endpoint for Prometheus
	s.router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
