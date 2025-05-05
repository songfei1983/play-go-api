package server

import (
	"context"
	"fmt"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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

	// Add CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		MaxAge:       86400, // 24小时
	}))

	e.Use(middleware.Logger())
	e.Use(otelecho.Middleware("api-service"))
	e.Use(middleware.RequestID())
	e.Use(mymiddleware.MetricsMiddleware())

	e.Logger.SetLevel(log.DEBUG)

	// Configure JWT middleware
	jwtConfig := echojwt.Config{
		SigningKey: []byte("your-secret-key"), // Replace with your secret key
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == "/health" ||
				c.Request().URL.Path == "/metrics" ||
				c.Request().URL.Path == "/api/v1/login" ||
				c.Request().URL.Path == "/api/v1/register"
		},
	}
	e.Use(echojwt.WithConfig(jwtConfig))

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

	// Register endpoints
	v1.POST("/register", userHandler.Register)
	v1.OPTIONS("/register", handleOptions)

	// Login endpoints
	v1.POST("/login", userHandler.Login)
	v1.OPTIONS("/login", handleOptions)

	// User management endpoints
	v1.GET("/users/current", userHandler.GetCurrentUser)
	v1.GET("/users/:id", userHandler.GetUser)
	v1.PUT("/users/:id", userHandler.UpdateUser)
	v1.DELETE("/users/:id", userHandler.SoftDeleteUser)
	v1.PATCH("/users/:id", userHandler.RestoreUser)
	v1.OPTIONS("/users/:id", handleOptions)

	// Product routes
	productHandler := handler.NewProductHandler(s.app.DB, s.app.Redis)
	products := v1.Group("/products")
	products.GET("", productHandler.List)
	products.POST("", productHandler.Create)
	products.GET("/:id", productHandler.Get)
	products.PUT("/:id", productHandler.Update)
	products.PATCH("/:id", productHandler.Update)
	products.DELETE("/:id/soft", productHandler.Delete)
	products.POST("/:id/restore", productHandler.Restore)

	// Metrics endpoint for Prometheus
	s.router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	s.router.OPTIONS("/metrics", handleOptions)

	// Health check endpoint (before JWT middleware)
	s.router.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"version": "1.0.0",
		})
	})
	s.router.OPTIONS("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}

// handleOptions handles OPTIONS requests for CORS
func handleOptions(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
