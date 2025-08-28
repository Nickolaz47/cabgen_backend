package testutils

import (
	"bytes"
	"encoding/json"
	"log"
	"maps"
	"net/http/httptest"

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

func CopyMap(original map[string]string) map[string]string {
	copy := make(map[string]string, len(original))
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
