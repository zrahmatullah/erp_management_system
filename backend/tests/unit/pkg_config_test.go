package unit_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"cafe-erp-system/backend/internal/config"
	"cafe-erp-system/backend/pkg/logger"
	"cafe-erp-system/backend/pkg/response"
)

func TestPkg_Response(t *testing.T) {
	t.Run("Success Response", func(t *testing.T) {
		rec := httptest.NewRecorder()
		data := map[string]string{"item": "espresso"}
		response.Success(rec, "Item retrieved", data)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var res response.Response
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assert.Equal(t, "Item retrieved", res.Message)
		assert.NotNil(t, res.Data)
	})

	t.Run("Error Response", func(t *testing.T) {
		rec := httptest.NewRecorder()
		response.Error(rec, http.StatusBadRequest, "Invalid input data")

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var res response.Response
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.False(t, res.Success)
		assert.Equal(t, "Invalid input data", res.Message)
	})

	t.Run("Paginated Response", func(t *testing.T) {
		rec := httptest.NewRecorder()
		items := []string{"Latte", "Cappuccino", "Mocha"}
		meta := response.PaginationMeta{
			Page:       1,
			PerPage:    10,
			Total:      3,
			TotalPages: 1,
		}
		response.Paginated(rec, "Products listed", items, meta)

		assert.Equal(t, http.StatusOK, rec.Code)
		var res response.Response
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assert.Equal(t, "Products listed", res.Message)
		assert.NotNil(t, res.Meta)
		assert.Equal(t, 1, res.Meta.Page)
		assert.Equal(t, 3, res.Meta.Total)
	})
}

func TestPkg_Logger(t *testing.T) {
	assert.NotPanics(t, func() {
		logger.InitLogger("development")
		logger.Log.Info().Msg("Test logger message in development")

		logger.InitLogger("production")
		logger.Log.Info().Msg("Test logger message in production")
	})
}

func TestConfig_LoadConfig(t *testing.T) {
	t.Run("Load with Defaults and Environment Overrides", func(t *testing.T) {
		os.Setenv("PORT", "9999")
		os.Setenv("ENV", "staging")
		os.Setenv("DATABASE_URL", "postgres://user:pass@host:5432/db?sslmode=require")
		os.Setenv("JWT_SECRET", "custom-secret-key-for-testing")
		os.Setenv("JWT_EXPIRE_HOURS", "48")
		defer func() {
			os.Unsetenv("PORT")
			os.Unsetenv("ENV")
			os.Unsetenv("DATABASE_URL")
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("JWT_EXPIRE_HOURS")
		}()

		cfg, err := config.LoadConfig()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "9999", cfg.Server.Port)
		assert.Equal(t, "staging", cfg.Server.Env)
		assert.Equal(t, "postgres://user:pass@host:5432/db?sslmode=require", cfg.DB.DSN)
		assert.Equal(t, "custom-secret-key-for-testing", cfg.JWT.Secret)
		assert.Equal(t, 48, cfg.JWT.Expire)
	})

	t.Run("Load with Synthesized DSN when DATABASE_URL is Empty", func(t *testing.T) {
		os.Unsetenv("DATABASE_URL")
		os.Setenv("DB_USERNAME", "myuser")
		os.Setenv("DB_PASSWORD", "mypass")
		os.Setenv("DB_HOST", "myhost")
		os.Setenv("DB_PORT", "5433")
		os.Setenv("DB_DATABASE", "mydb")
		defer func() {
			os.Unsetenv("DB_USERNAME")
			os.Unsetenv("DB_PASSWORD")
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("DB_DATABASE")
		}()

		cfg, err := config.LoadConfig()
		assert.NoError(t, err)
		assert.Contains(t, cfg.DB.DSN, "postgres://myuser:mypass@myhost:5433/mydb")
	})
}
