package handler

import (
	"encoding/json"
	"errors"
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

// setupTest 初始化测试环境
func setupTest(t *testing.T) (*echo.Echo, *UserHandler, sqlmock.Sqlmock, redismock.ClientMock) {
	// 初始化 Echo 框架
	e := echo.New()

	// 初始化 SQL Mock
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	// 使用 GORM 连接 SQL Mock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	// 初始化 Redis Mock
	redisMock, redisMockClient := redismock.NewClientMock()

	// 创建 UserHandler 实例
	// 将 redismock.ClientMock 转换为 *redis.Client
	// 将 redismock.ClientMock 转换为 redis.Client 指针
	userHandler := NewUserHandler(gormDB, redisMock)

	return e, userHandler, mock, redisMockClient
}

// TestRegister 测试用户注册功能
func TestRegister(t *testing.T) {
	// 设置测试环境
	e, handler, mock, _ := setupTest(t)

	t.Run("成功注册用户", func(t *testing.T) {
		userJSON := `{"username":"testuser","password":"password123","email":"test@example.com"}`

		// 创建请求
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// 设置数据库期望
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// 执行请求
		err := handler.Register(c)

		// 断言结果
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// 验证响应内容
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(1), response["id"])
		assert.Equal(t, "testuser", response["username"])
		assert.Equal(t, "test@example.com", response["email"])
	})

	t.Run("缺少必填字段", func(t *testing.T) {
		// 准备请求数据（缺少必填字段）
		userJSON := `{"username":"testuser"}`

		// 创建请求
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// 执行请求
		err := handler.Register(c)

		// 断言结果
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// 验证响应内容
		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "required")
	})

	t.Run("数据库错误", func(t *testing.T) {
		// 准备请求数据
		userJSON := `{"username":"testuser","password":"password123","email":"test@example.com"}`

		// 创建请求
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// 设置数据库期望（返回错误）
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `users`").WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		// 执行请求
		err := handler.Register(c)

		// 断言结果
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// 验证响应内容
		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "database error")
	})
}

// TestGetUser 测试获取单个用户功能
func TestGetUser(t *testing.T) {
	// 设置测试环境
	e, handler, mock, redisMock := setupTest(t)

	t.Run("从缓存获取用户", func(t *testing.T) {
		// 准备用户数据
		user := User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Status:   "active",
		}
		userJSON, _ := json.Marshal(user)

		// 创建请求
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// 设置Redis期望
		redisMock.ExpectGet("users:1").SetVal(string(userJSON))

		// 执行请求
		err := handler.Get(c)

		// 断言结果
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// 验证响应内容
		var response User
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "testuser", response.Username)
		assert.Equal(t, "test@example.com", response.Email)
	})

	t.Run("从数据库获取用户", func(t *testing.T) {
		// 创建请求
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// 设置Redis期望（返回nil，表示缓存未命中）
		redisMock.ExpectGet("user:1").RedisNil()

		// 设置数据库期望
		rows := sqlmock.NewRows([]string{"id", "username", "password", "email", "first_name", "last_name", "phone", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "testuser", "password123", "test@example.com", "Test", "User", "1234567890", "active", time.Now(), time.Now(), nil)

		mock.ExpectQuery("SELECT \\* FROM `users` WHERE deleted_at IS NULL AND `users`\\.`id` = \\? ORDER BY `users`\\.`id` LIMIT \\?").
			WithArgs(1, 1).
			WillReturnRows(rows)

		// 设置Redis Set期望（缓存用户数据）
		redisMock.ExpectSet("user:1", sqlmock.AnyArg(), time.Hour).SetVal("OK")

		// 执行请求
		err := handler.GetUser(c)

		// 断言结果
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// 验证响应内容
		var response User
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "testuser", response.Username)
	})

	t.Run("用户不存在", func(t *testing.T) {
		// 创建请求
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("999")

		// 设置Redis期望（返回nil，表示缓存未命中）
		redisMock.ExpectGet("user:999").RedisNil()

		// 设置数据库期望（返回错误，表示用户不存在）
		mock.ExpectQuery("SELECT \\* FROM `users` WHERE deleted_at IS NULL AND `users`\\.`id` = \\? ORDER BY `users`\\.`id` LIMIT \\?").
			WithArgs(999, 1).
			WillReturnRows(sqlmock.NewRows([]string{}))

		// 执行请求
		err := handler.GetUser(c)

		// 断言结果
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		// 验证响应内容
		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "not found")
	})
}

// TestGetUsers 测试获取用户列表功能
func TestGetUsers(t *testing.T) {
	// 设置测试环境
	e, handler, mock, _ := setupTest(t)

	t.Run("获取所有用户", func(t *testing.T) {
		// 创建请求
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// 设置数据库期望
		rows := sqlmock.NewRows([]string{"id", "username", "password", "email", "first_name", "last_name", "phone", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "user1", "password1", "user1@example.com", "First1", "Last1", "1111111111", "active", time.Now(), time.Now(), nil).
			AddRow(2, "user2", "password2", "user2@example.com", "First2", "Last2", "2222222222", "active", time.Now(), time.Now(), nil)

		mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE deleted_at IS NULL$").WillReturnRows(rows)

		// 执行请求
		err := handler.GetUsers(c)

		// 断言结果
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// 验证响应内容
		var response []User
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, uint(1), response[0].ID)
		assert.Equal(t, uint(2), response[1].ID)
	})

	t.Run("包含已删除用户", func(t *testing.T) {
		// 创建请求
		req := httptest.NewRequest(http.MethodGet, "/users?include_deleted=true", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.QueryParams().Set("include_deleted", "true")

		// 设置数据库期望
		deleteTime := time.Now()
		rows := sqlmock.NewRows([]string{"id", "username", "password", "email", "first_name", "last_name", "phone", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "user1", "password1", "user1@example.com", "First1", "Last1", "1111111111", "active", time.Now(), time.Now(), nil).
			AddRow(2, "user2", "password2", "user2@example.com", "First2", "Last2", "2222222222", "active", time.Now(), time.Now(), nil).
			AddRow(3, "user3", "password3", "user3@example.com", "First3", "Last3", "3333333333", "inactive", time.Now(), time.Now(), &deleteTime)

		mock.ExpectQuery("^SELECT (.+) FROM `users`$").WillReturnRows(rows)

		// 执行请求
		err := handler.GetUsers(c)

		// 断言结果
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// 验证响应内容
		var response []User
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 3)
		assert.Equal(t, uint(3), response[2].ID)
		assert.NotNil(t, response[2].DeletedAt)
	})
}

// TestUpdateUser 测试更新用户功能
func TestUpdateUser(t *testing.T) {
	// 设置测试环境
	e, handler, mock, redisMock := setupTest(t)

	t.Run("成功更新用户", func(t *testing.T) {
		// 准备请求数据
		updateJSON := `{"first_name":"Updated","last_name":"Name"}`

		// 创建请求
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(updateJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// 设置数据库期望
		rows := sqlmock.NewRows([]string{"id", "username", "password", "email", "first_name", "last_name", "phone", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "testuser", "password123", "test@example.com", "Test", "User", "1234567890", "active", time.Now(), time.Now(), nil)

		mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`\\.`id` = \\? ORDER BY `users`\\.`id` LIMIT \\?").
			WithArgs(1, 1).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `users` SET").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// 设置Redis期望（删除缓存）
		redisMock.ExpectDel("user:1").SetVal(1)

		// 执行请求
		err := handler.UpdateUser(c)

		// 断言结果
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("用户不存在", func(t *testing.T) {
		// 准备请求数据
		updateJSON := `{"first_name":"Updated","last_name":"Name"}`

		// 创建请求
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(updateJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("999")

		// 设置数据库期望（返回错误，表示用户不存在）
		mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE (.+) LIMIT 1$").WillReturnError(gorm.ErrRecordNotFound)

		// 执行请求
		err := handler.UpdateUser(c)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)
	})
}

// TestSoftDeleteUser 测试软删除用户功能
func TestSoftDeleteUser(t *testing.T) {
	// 设置测试环境
	e, handler, mock, redisMock := setupTest(t)

	t.Run("成功软删除用户", func(t *testing.T) {
		// 创建请求
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// 设置数据库期望
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `users` SET").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// 设置Redis期望（删除缓存）
		redisMock.ExpectDel("user:1").SetVal(1)

		// 执行请求
		err := handler.SoftDeleteUser(c)

		// 断言结果
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("删除不存在的用户", func(t *testing.T) {
		// 创建请求
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("999")

		// 设置数据库期望（返回错误，表示用户不存在）
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE `users` SET (.+) WHERE (.+)$").WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		// 执行请求
		err := handler.SoftDeleteUser(c)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)
	})
}

// TestRestoreUser 测试恢复已删除用户功能
func TestRestoreUser(t *testing.T) {
	// 设置测试环境
	e, handler, mock, _ := setupTest(t)

	t.Run("成功恢复用户", func(t *testing.T) {
		// 创建请求
		req := httptest.NewRequest(http.MethodPatch, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// 设置数据库期望
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE `users` SET (.+) WHERE (.+)$").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// 执行请求
		err := handler.RestoreUser(c)

		// 断言结果
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("恢复不存在的用户", func(t *testing.T) {
		// 创建请求
		req := httptest.NewRequest(http.MethodPatch, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("999")

		// 设置数据库期望（返回错误，表示用户不存在）
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE `users` SET (.+) WHERE (.+)$").WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		// 执行请求
		err := handler.RestoreUser(c)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	})
}

// TestLogin 测试用户登录功能
func TestLogin(t *testing.T) {
	// 设置测试环境
	e, handler, mock, _ := setupTest(t)

	t.Run("登录成功", func(t *testing.T) {
		loginJSON := `{"username":"testuser","password":"password123"}`

		// 创建请求
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// 设置数据库期望
		rows := sqlmock.NewRows([]string{"id", "username", "password", "email", "status"}).
			AddRow(1, "testuser", "password123", "test@example.com", "active")

		mock.ExpectQuery("SELECT \\* FROM `users` WHERE username = \\? ORDER BY `users`\\.`id` LIMIT \\?").
			WithArgs("testuser", 1).
			WillReturnRows(rows)

		// 执行请求
		err := handler.Login(c)

		// 断言结果
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// 验证响应内容
		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "token")
		assert.NotEmpty(t, response["token"])
	})

	t.Run("用户不存在", func(t *testing.T) {
		// 准备请求数据
		loginJSON := `{"username":"nonexistent","password":"password123"}`

		// 创建请求
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// 设置数据库期望（返回错误，表示用户不存在）
		mock.ExpectQuery("SELECT \\* FROM `users` WHERE username = \\? ORDER BY `users`\\.`id` LIMIT \\?").
			WithArgs("nonexistent", 1).
			WillReturnRows(sqlmock.NewRows([]string{}))

		// 执行请求
		err := handler.Login(c)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.(*echo.HTTPError).Code)
		assert.Contains(t, err.(*echo.HTTPError).Message, "invalid credentials")
	})

	t.Run("密码错误", func(t *testing.T) {
		// 准备请求数据
		loginJSON := `{"username":"testuser","password":"wrongpassword"}`

		// 创建请求
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// 设置数据库期望
		rows := sqlmock.NewRows([]string{"id", "username", "password", "email", "first_name", "last_name", "phone", "status", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "testuser", "password123", "test@example.com", "Test", "User", "1234567890", "active", time.Now(), time.Now(), nil)

		mock.ExpectQuery("SELECT \\* FROM `users` WHERE username = \\? ORDER BY `users`\\.`id` LIMIT \\?").
			WithArgs("testuser", 1).
			WillReturnRows(rows)

		// 执行请求
		err := handler.Login(c)

		// 断言结果
		assert.Error(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.(*echo.HTTPError).Code)
		assert.Contains(t, err.(*echo.HTTPError).Message, "invalid credentials")
	})
}
