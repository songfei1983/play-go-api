package handler

import (
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Product 产品模型
type Product struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	Price       float64    `json:"price" gorm:"not null"`
	Stock       int        `json:"stock" gorm:"not null"`
	Status      string     `json:"status" gorm:"default:active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// GetID 实现 Model 接口
func (p Product) GetID() uint {
	return p.ID
}

// TableName 实现 Model 接口
func (p Product) TableName() string {
	return "products"
}

// ProductHandler 产品处理器
type ProductHandler struct {
	*BaseHandler[Product]
}

// NewProductHandler 创建产品处理器
func NewProductHandler(db *gorm.DB, redis *redis.Client) *ProductHandler {
	return &ProductHandler{
		BaseHandler: NewBaseHandler[Product](db, redis),
	}
}
