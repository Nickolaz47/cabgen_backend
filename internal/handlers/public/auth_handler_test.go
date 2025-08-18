package public_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	testutils.SetupTestContext()

	db := testutils.SetupTestRepos()

	db.Create(&models.Country{
		Code: "BRA",
		Pt:   "Brasil",
		En:   "Brazil",
		Es:   "Brazil",
	})
	t.Run("Success", func(t *testing.T) {
		body := `{"name":"Nicolas","username":"nmfaraujo","email":"nicolas@mail.com","confirm_email":"nicolas@mail.com","password":"12345678","confirm_password":"12345678","country_code":"BRA"}`
		expected := map[string]any{
			"message": "User created successfully. Please wait for activation by an administrator.",
			"data": map[string]any{
				"name":         "Nicolas",
				"username":     "nmfaraujo",
				"email":        "nicolas@mail.com",
				"country_code": "BRA",
				"country":      "Brazil",
				"user_role":    "Collaborator",
			},
		}

		c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/register", body)
		public.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var got map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)

		if data, ok := got["data"].(map[string]any); ok {
			delete(data, "id")
		}

		assert.Equal(t, expected, got)
	})

	for _, tt := range data.RegisterTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/register", tt.Body)
			public.Register(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Email already exists", func(t *testing.T) {
		rawBody := map[string]string{
			"name":             "Nicolas",
			"username":         "nick",
			"email":            "nicolas@mail.com",
			"confirm_email":    "nicolas@mail.com",
			"password":         "12345678",
			"confirm_password": "12345678",
			"country_code":     "BRA",
			"user_role":        "Collaborator",
		}
		body := testutils.ToJSON(rawBody)
		expected := `{"error": "Email is already in use."}`

		c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/register", body)
		public.Register(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Username already exists", func(t *testing.T) {
		rawBody := map[string]string{
			"name":             "Nicolas",
			"username":         "nmfaraujo",
			"email":            "nick@mail.com",
			"confirm_email":    "nick@mail.com",
			"password":         "12345678",
			"confirm_password": "12345678",
			"country_code":     "BRA",
			"user_role":        "Collaborator",
		}
		body := testutils.ToJSON(rawBody)
		expected := `{"error": "Username already exists."}`

		c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/register", body)
		public.Register(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
func TestLogin(t *testing.T) {
	testutils.SetupTestContext()

	db := testutils.SetupTestRepos()

	db.Create(&models.User{
		Name:        "Nicolas",
		Username:    "nick",
		Password:    "$2a$10$P8SRTHBxlK09pYuj8Nn1A.2WMufAH1tZZKAPQel1bt0X5S82zbRGO",
		Email:       "nick@mail.com",
		CountryCode: "BRA",
		Country:     models.Country{Code: "BRA", Pt: "Brasil", Es: "Brazil", En: "Brazil"},
		IsActive:    true,
	})

	t.Run(data.LoginSuccess.Name, func(t *testing.T) {
		c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/login", data.LoginSuccess.Body)
		public.Login(c)

		cookies := w.Result().Cookies()
		var accessCookie, refreshCookie string

		for _, cookie := range cookies {
			if cookie.Name == auth.Access {
				accessCookie = cookie.Value
			}
			if cookie.Name == auth.Refresh {
				refreshCookie = cookie.Value
			}
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, data.LoginSuccess.Expected, w.Body.String())
		assert.NotEmpty(t, accessCookie)
		assert.NotEmpty(t, refreshCookie)
	})

	for _, tt := range data.LoginBadRequestTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/login", tt.Body)
			public.Login(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	for _, tt := range data.LoginUnauthorizedTests {
		t.Run(tt.Name, func(t *testing.T) {
			c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/login", tt.Body)
			public.Login(c)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	db.Model(&models.User{}).Where("username = ?", "nick").Update("is_active", false)

	t.Run(data.LoginUserDeactivatedTest.Name, func(t *testing.T) {
		c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/login", data.LoginUserDeactivatedTest.Body)
		public.Login(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.JSONEq(t, data.LoginUserDeactivatedTest.Expected, w.Body.String())
	})
}
func TestLogout(t *testing.T)  {}
func TestRefresh(t *testing.T) {}
