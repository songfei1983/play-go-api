package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v8"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupProductTest(t *testing.T) (*echo.Echo, *ProductHandler, sqlmock.Sqlmock, redismock.ClientMock) {
	// Initialize Echo framework
	e := echo.New()

	// Initialize SQL Mock
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	// Use GORM with SQL Mock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	// Initialize Redis Mock
	redisMock, redisMockClient := redismock.NewClientMock()

	// Create ProductHandler instance
	productHandler := NewProductHandler(gormDB, redisMock)

	return e, productHandler, mock, redisMockClient
}

func TestProductCRUD(t *testing.T) {
	e, handler, mock, redisMock := setupProductTest(t)

	t.Run("创建产品", func(t *testing.T) {
		productJSON := `{
			"name": "Test Product",
			"description": "Test Description",
			"price": 99.99,
			"stock": 100,
			"status": "active"
		}`

		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(productJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `products`").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := handler.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response Product
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "Test Product", response.Name)
	})

	t.Run("获取产品列表", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "Product 1", "Description 1", 99.99, 100, "active", time.Now(), time.Now(), nil).
			AddRow(2, "Product 2", "Description 2", 199.99, 50, "active", time.Now(), time.Now(), nil)

		mock.ExpectQuery("^SELECT (.+) FROM `products` WHERE deleted_at IS NULL").WillReturnRows(rows)

		err := handler.List(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response []Product
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
	})

	t.Run("获取单个产品", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/products/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// Test cache miss scenario
		redisMock.ExpectGet("products:1").RedisNil()

		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "Test Product", "Test Description", 99.99, 100, "active", time.Now(), time.Now(), nil)

		mock.ExpectQuery("SELECT \\* FROM `products` WHERE deleted_at IS NULL AND `products`\\.`id` = \\? ORDER BY `products`\\.`id` LIMIT \\?").
			WithArgs(1, 1).
			WillReturnRows(rows)

		// Expect cache set
		redisMock.ExpectSet("products:1", sqlmock.AnyArg(), time.Hour).SetVal("OK")

		err := handler.Get(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response Product
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "Test Product", response.Name)
	})

	t.Run("更新产品", func(t *testing.T) {
		productJSON := `{
			"name": "Updated Product",
			"description": "Updated Description",
			"price": 149.99,
			"stock": 75
		}`

		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(productJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/products/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock", "status"}).
			AddRow(1, "Test Product", "Test Description", 99.99, 100, "active")

		mock.ExpectQuery("SELECT \\* FROM `products` WHERE `products`\\.`id` = \\? ORDER BY `products`\\.`id` LIMIT \\?").
			WithArgs(1, 1).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `products`").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Expect cache deletion
		redisMock.ExpectDel("products:1").SetVal(1)

		err := handler.Update(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("软删除产品", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/products/:id/soft")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `products` SET `deleted_at`=(.+),`updated_at`=(.+) WHERE id = \\?").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Expect cache deletion
		redisMock.ExpectDel("products:1").SetVal(1)

		err := handler.Delete(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("恢复已删除产品", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/products/:id/restore")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `products` SET `deleted_at`=NULL,`updated_at`=(.+) WHERE id = \\?").
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := handler.Restore(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
