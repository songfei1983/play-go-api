package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewUserHandler(db *gorm.DB, redis *redis.Client) *UserHandler {
	return &UserHandler{
		db:    db,
		redis: redis,
	}
}

func (h *UserHandler) Register(c echo.Context) error {
	ctx := c.Request().Context()
	tracer := otel.Tracer("api-service")
	ctx, span := tracer.Start(ctx, "UserHandler.Register")
	span.SetAttributes(attribute.String("Request-ID", c.Response().Header().Get("X-Request-ID")))
	defer span.End()

	var user User
	if err := c.Bind(&user); err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.db.WithContext(ctx).Create(&user).Error; err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	ctx := c.Request().Context()
	tracer := otel.Tracer("api-service")
	ctx, span := tracer.Start(ctx, "UserHandler.GetUser")
	span.SetAttributes(attribute.String("Request-ID", c.Response().Header().Get("X-Request-ID")))
	defer span.End()

	// random delay
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	id := c.Param("id")
	span.SetAttributes(attribute.String("user_id", id))

	// Try to get from cache first
	cacheKey := fmt.Sprintf("user:%s", id)
	userJSON, err := h.redis.Get(ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		span.RecordError(err)
		span.SetAttributes(semconv.HTTPResponseStatusCode(http.StatusInternalServerError))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Cache error"})
	}
	if err == nil {
		var user User
		if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
			span.RecordError(err)
			span.SetAttributes(semconv.HTTPResponseStatusCode(http.StatusInternalServerError))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Cache unmarshal error"})
		}
		span.SetAttributes(semconv.HTTPResponseStatusCode(http.StatusOK))
		return c.JSON(http.StatusOK, user)
	}

	var user User
	if err := h.db.WithContext(ctx).First(&user, id).Error; err != nil {
		span.RecordError(err)
		span.SetAttributes(semconv.HTTPResponseStatusCode(http.StatusNotFound))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Cache the user
	userByte, err := json.Marshal(user)
	if err != nil {
		span.RecordError(err)
	} else {
		if err := h.redis.Set(ctx, cacheKey, string(userByte), time.Hour).Err(); err != nil {
			span.RecordError(err)
		}
	}

	span.SetAttributes(semconv.HTTPResponseStatusCode(http.StatusOK))
	return c.JSON(http.StatusOK, user)
}
