package public_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

var MockRegisterUser = models.RegisterInput{
	Name:            "Nicolas",
	Username:        "nmfaraujo",
	Password:        "12345678",
	ConfirmPassword: "12345678",
	Email:           "nicolas@mail.com",
	ConfirmEmail:    "nicolas@mail.com",
	CountryCode:     "BRA",
}

var MockLoginUser = models.User{
	Name:        "Nicolas",
	Username:    "nick",
	Password:    "$2a$10$P8SRTHBxlK09pYuj8Nn1A.2WMufAH1tZZKAPQel1bt0X5S82zbRGO",
	Email:       "nick@mail.com",
	CountryCode: "BRA",
	Country:     models.Country{Code: "BRA", Pt: "Brasil", Es: "Brazil", En: "Brazil"},
	IsActive:    true,
}

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
		body := testutils.ToJSON(MockRegisterUser)
		expected := map[string]any{
			"message": "User created successfully. Please wait for activation by an administrator.",
			"data": map[string]any{
				"name":         MockRegisterUser.Name,
				"username":     MockRegisterUser.Username,
				"email":        MockRegisterUser.Email,
				"country_code": MockRegisterUser.CountryCode,
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
		MockRegisterUser.Username = "nick"
		body := testutils.ToJSON(MockRegisterUser)
		expected := `{"error": "Email is already in use."}`

		c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/register", body)
		public.Register(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Username already exists", func(t *testing.T) {
		MockRegisterUser.Email = "nick@mail.com"
		MockRegisterUser.ConfirmEmail = "nick@mail.com"
		MockRegisterUser.Username = "nmfaraujo"
		body := testutils.ToJSON(MockRegisterUser)
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

	db.Create(&MockLoginUser)

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

	db.Model(&models.User{}).
		Where("username = ?", MockLoginUser.Username).
		Update("is_active", false)

	t.Run(data.LoginUserDeactivatedTest.Name, func(t *testing.T) {
		c, w := testutils.SetupGinContext(http.MethodPost, "/api/auth/login", data.LoginUserDeactivatedTest.Body)
		public.Login(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.JSONEq(t, data.LoginUserDeactivatedTest.Expected, w.Body.String())
	})
}
func TestLogout(t *testing.T) {
	db := testutils.SetupTestRepos()

	db.Create(&MockLoginUser)
	db.Model(&models.User{}).
		Where("username = ?", MockLoginUser.Username).
		Update("is_active", true)

	c, w := testutils.SetupGinContext(
		http.MethodPost,
		"/api/auth/login",
		testutils.ToJSON(models.LoginInput{
			Username: MockLoginUser.Username,
			Password: MockRegisterUser.Password},
		),
	)

	public.Login(c)
	var accessCookie, refreshCookie string

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == auth.Access {
			accessCookie = cookie.Value
		}
		if cookie.Name == auth.Refresh {
			refreshCookie = cookie.Value
		}
	}

	assert.NotEmpty(t, accessCookie)
	assert.NotEmpty(t, refreshCookie)

	body := bytes.NewBuffer(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", body)
	req.AddCookie(&http.Cookie{Name: auth.Access, Value: accessCookie})
	req.AddCookie(&http.Cookie{Name: auth.Refresh, Value: refreshCookie})

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = req
	
	public.Logout(c)

	cookies = w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == auth.Access {
			accessCookie = cookie.Value
		}
		if cookie.Name == auth.Refresh {
			refreshCookie = cookie.Value
		}
	}

	expected := `{"message": "Logout successful."}`

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
	assert.Empty(t, accessCookie)
	assert.Empty(t, refreshCookie)
}
func TestRefresh(t *testing.T) {}
