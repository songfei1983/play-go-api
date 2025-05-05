package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"
)

// Model 定义基础模型接口
type Model interface {
	GetID() uint
	TableName() string
}

// BaseHandler 通用CRUD处理器
type BaseHandler[T Model] struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewBaseHandler 创建基础处理器
func NewBaseHandler[T Model](db *gorm.DB, redis *redis.Client) *BaseHandler[T] {
	return &BaseHandler[T]{
		db:    db,
		redis: redis,
	}
}

// getCacheKey 获取缓存键
func (h *BaseHandler[T]) getCacheKey(id string) string {
	var model T
	return fmt.Sprintf("%s:%s", model.TableName(), id)
}

// Create 通用创建方法
func (h *BaseHandler[T]) Create(c echo.Context) error {
	ctx := c.Request().Context()
	tracer := otel.Tracer("api-service")
	ctx, span := tracer.Start(ctx, "BaseHandler.Create")
	span.SetAttributes(attribute.String("Request-ID", c.Response().Header().Get("X-Request-ID")))
	defer span.End()

	var model T
	if err := c.Bind(&model); err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	result := h.db.WithContext(ctx).Create(&model)
	if result.Error != nil {
		span.RecordError(result.Error)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}

	return c.JSON(http.StatusCreated, model)
}

// Get 通用获取单个记录方法
func (h *BaseHandler[T]) Get(c echo.Context) error {
	ctx := c.Request().Context()
	tracer := otel.Tracer("api-service")
	ctx, span := tracer.Start(ctx, "BaseHandler.Get")
	span.SetAttributes(attribute.String("Request-ID", c.Response().Header().Get("X-Request-ID")))
	defer span.End()

	id := c.Param("id")
	span.SetAttributes(attribute.String("id", id))

	// 尝试从缓存获取
	cacheKey := h.getCacheKey(id)
	modelJSON, err := h.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var model T
		if err = json.Unmarshal([]byte(modelJSON), &model); err == nil {
			span.SetAttributes(attribute.String("data_source", "cache"))
			return c.JSON(http.StatusOK, model)
		}
		span.RecordError(err)
	}

	var model T
	intID, err := strconv.Atoi(id)
	if err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	if err = h.db.WithContext(ctx).Where("deleted_at IS NULL").First(&model, intID).Error; err != nil {
		span.RecordError(err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Record not found"})
	}

	span.SetAttributes(attribute.String("data_source", "mysql"))

	// 缓存结果
	if modelByte, err := json.Marshal(model); err == nil {
		h.redis.Set(ctx, cacheKey, string(modelByte), time.Hour)
	}

	return c.JSON(http.StatusOK, model)
}

// List 通用获取列表方法
func (h *BaseHandler[T]) List(c echo.Context) error {
	ctx := c.Request().Context()
	tracer := otel.Tracer("api-service")
	ctx, span := tracer.Start(ctx, "BaseHandler.List")
	defer span.End()

	var models []T
	query := h.db.WithContext(ctx)

	includeSoftDeleted := c.QueryParam("include_deleted") == "true"
	if !includeSoftDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	if err := query.Find(&models).Error; err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, models)
}

// Update 通用更新方法
func (h *BaseHandler[T]) Update(c echo.Context) error {
	ctx := c.Request().Context()
	tracer := otel.Tracer("api-service")
	ctx, span := tracer.Start(ctx, "BaseHandler.Update")
	defer span.End()

	id := c.Param("id")
	var model T

	intID, err := strconv.Atoi(id)
	if err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	if err := h.db.WithContext(ctx).First(&model, intID).Error; err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusNotFound, "Record not found")
	}

	if c.Request().Method == "PUT" {
		if err := c.Bind(&model); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	} else {
		updates := make(map[string]interface{})
		if err := c.Bind(&updates); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := h.db.Model(&model).Updates(updates).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// 清除缓存
	cacheKey := h.getCacheKey(id)
	h.redis.Del(ctx, cacheKey)

	return c.JSON(http.StatusOK, model)
}

// Delete 通用软删除方法
func (h *BaseHandler[T]) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	tracer := otel.Tracer("api-service")
	ctx, span := tracer.Start(ctx, "BaseHandler.Delete")
	defer span.End()

	id := c.Param("id")

	intID, err := strconv.Atoi(id)
	if err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	var model T
	if err := h.db.WithContext(ctx).Model(&model).Unscoped().Where("id = ?", intID).Update("deleted_at", time.Now()).Error; err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// 清除缓存
	cacheKey := h.getCacheKey(id)
	h.redis.Del(ctx, cacheKey)

	return c.NoContent(http.StatusNoContent)
}

// Restore 通用恢复方法
func (h *BaseHandler[T]) Restore(c echo.Context) error {
	ctx := c.Request().Context()
	tracer := otel.Tracer("api-service")
	ctx, span := tracer.Start(ctx, "BaseHandler.Restore")
	defer span.End()

	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	var model T
	if err := h.db.WithContext(ctx).Model(&model).Unscoped().Where("id = ?", intID).Update("deleted_at", nil).Error; err != nil {
		span.RecordError(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
