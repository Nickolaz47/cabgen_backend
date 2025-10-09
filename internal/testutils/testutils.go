package testutils

import (
	"bytes"
	"encoding/json"
	"log"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewMockDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Exec("PRAGMA foreign_keys = ON")
	db.AutoMigrate(&models.Country{})
	db.AutoMigrate(&testmodels.User{})
	db.AutoMigrate(&testmodels.Origin{})
	db.AutoMigrate(&testmodels.Sequencer{})

	return db
}

func SetupTestRepos() *gorm.DB {
	db := NewMockDB()
	repository.InitRepositories(db)
	return db
}

func SetupGinContext(method, URL, body string, headers map[string]string, params gin.Params) (*gin.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(
		method,
		URL,
		bytes.NewBufferString(body),
	)

	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = params

	return c, w
}

func SetupTestContext() {
	gin.SetMode(gin.TestMode)
	translation.LoadTranslation()
}

func CopyMap(original map[string]any) map[string]any {
	copy := make(map[string]any, len(original))
	maps.Copy(copy, original)
	return copy
}

func ToJSON(body any) string {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Failed to convert body to JSON: %v", err)
	}
	return string(jsonBytes)
}

func SetupMiddlewareContext() (*httptest.ResponseRecorder, *gin.Engine) {
	return httptest.NewRecorder(), gin.New()
}

func AddMiddlewares(engine *gin.Engine, middlewares ...gin.HandlerFunc) {
	for _, m := range middlewares {
		engine.Use(m)
	}
}

func AddTestGetRoute(engine *gin.Engine, statusCode int) {
	engine.GET("/", func(c *gin.Context) {
		if statusCode <= 0 {
			c.Status(http.StatusOK)
		} else {
			c.Status(statusCode)
		}
	})
}

func DoGetRequest(r *gin.Engine, w *httptest.ResponseRecorder) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
}

func WriteMockEnvFile(t *testing.T, envFilePath, envContent string) {
	if err := os.WriteFile(envFilePath, []byte(envContent), 0644); err != nil {
		t.Errorf("failed to write mock env file: %v", err)
	}
}
