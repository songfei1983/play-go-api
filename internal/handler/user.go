package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5" // Replace dgrijalva/jwt-go with this
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"gorm.io/gorm"
)

type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Username  string     `json:"username" gorm:"unique"`
	Password  string     `json:"password" gorm:"not null"` // "-" to exclude from JSON
	Email     string     `json:"email" gorm:"unique"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Phone     string     `json:"phone"`
	Status    string     `json:"status" gorm:"default:active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate required fields
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username, password and email are required"})
	}

	// Create user
	result := h.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		span.RecordError(result.Error)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}

	// Return created user with ID
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"status":   user.Status,
	})
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

func (h *UserHandler) GetUsers(c echo.Context) error {
	ctx := c.Request().Context()
	tracer := otel.Tracer("api-service")
	ctx, span := tracer.Start(ctx, "UserHandler.GetUsers")
	defer span.End()

	var users []User
	query := h.db.WithContext(ctx)

	includeSoftDeleted := c.QueryParam("include_deleted") == "true"
	if !includeSoftDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	if err := query.Find(&users).Error; err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	tracer := otel.Tracer("api-service")
	ctx, span := tracer.Start(ctx, "UserHandler.UpdateUser")
	defer span.End()

	id := c.Param("id")
	var user User

	if err := h.db.WithContext(ctx).First(&user, id).Error; err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// Handle both PUT (full update) and PATCH (partial update)
	if c.Request().Method == "PUT" {
		if err := c.Bind(&user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	} else {
		// For PATCH, only update provided fields
		updates := make(map[string]interface{})
		if err := c.Bind(&updates); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := h.db.Model(&user).Updates(updates).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// Clear cache
	cacheKey := fmt.Sprintf("user:%s", id)
	h.redis.Del(ctx, cacheKey)

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) SoftDeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	tracer := otel.Tracer("api-service")
	ctx, span := tracer.Start(ctx, "UserHandler.SoftDeleteUser")
	defer span.End()

	id := c.Param("id")

	if err := h.db.WithContext(ctx).Delete(&User{}, id).Error; err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Clear cache
	cacheKey := fmt.Sprintf("user:%s", id)
	h.redis.Del(ctx, cacheKey)

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) RestoreUser(c echo.Context) error {
	ctx := c.Request().Context()
	tracer := otel.Tracer("api-service")
	ctx, span := tracer.Start(ctx, "UserHandler.RestoreUser")
	defer span.End()

	id := c.Param("id")

	if err := h.db.WithContext(ctx).Model(&User{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	tracer := otel.Tracer("api-service")
	ctx, span := tracer.Start(ctx, "UserHandler.Login")
	span.SetAttributes(attribute.String("Request-ID", c.Response().Header().Get("X-Request-ID")))
	defer span.End()

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var user User
	if err := h.db.WithContext(ctx).Where("username = ?", req.Username).First(&user).Error; err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	if user.Password != req.Password {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	// Create claims with user data
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate signed token
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not generate token")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})
}
