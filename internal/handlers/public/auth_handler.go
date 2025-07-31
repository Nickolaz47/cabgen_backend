package public

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/auth"
	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var newUser models.RegisterInput
	if errMsg, valid := validations.ValidateRegisterInput(c, localizer, &newUser); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	var existingUser models.User
	err := db.DB.Where("email = ? OR username = ?", newUser.Email, newUser.Username).
		First(&existingUser).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	if err == nil {
		var errMsg string
		if existingUser.Email == newUser.Email {
			errMsg = responses.RegisterEmailAlreadyExistsError
		} else if existingUser.Username == newUser.Username {
			errMsg = responses.RegisterUsernameAlreadyExistsError
		}

		c.JSON(http.StatusConflict,
			responses.APIResponse{Error: responses.GetResponse(localizer, errMsg)},
		)
		return
	}

	if newUser.Email != newUser.ConfirmEmail {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.RegisterValidationEmailMismatch)},
		)
		return
	}

	if newUser.Password != newUser.ConfirmPassword {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.RegisterValidationPasswordMismatch)},
		)
		return
	}

	hashedPassword, err := security.Hash(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	user := models.User{
		Name:        newUser.Name,
		Username:    newUser.Username,
		Email:       newUser.Email,
		Password:    hashedPassword,
		CountryCode: newUser.CountryCode,
		IsActive:    false,
		UserRole:    models.Collaborator,
		Interest:    newUser.Interest,
		Role:        newUser.Role,
		Institution: newUser.Institution,
		CreatedBy:   newUser.Username,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.RegisterCreateUserError)},
		)
		return
	}

	c.JSON(http.StatusCreated,
		responses.APIResponse{
			Data:    user.ToResponse(),
			Message: responses.GetResponse(localizer, responses.RegisterMessage),
		},
	)
}

func Login(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var login models.LoginInput
	if errMsg, valid := validations.ValidateLoginInput(c, localizer, &login); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	var existingUser models.User
	err := db.DB.Where("username = ?", login.Username).First(&existingUser).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.LoginInvalidCredentialsError)})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	if err = security.CheckPassword(existingUser.Password, login.Password); err != nil {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.LoginInvalidCredentialsError)})
		return
	}

	accessKey, err := auth.GetSecretKey(auth.Access)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
		return
	}

	refreshKey, err := auth.GetSecretKey(auth.Refresh)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
		return
	}

	accessToken, err := auth.GenerateToken(existingUser.ToToken(), accessKey, auth.AccessTokenExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
		return
	}

	refreshToken, err := auth.GenerateToken(existingUser.ToToken(), refreshKey, auth.RefreshTokenExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
		return
	}

	accessCookie := auth.CreateCookie(auth.Access, accessToken, "/", auth.AccessTokenExpiration)
	refreshCookie := auth.CreateCookie(auth.Refresh, refreshToken, "/api/auth/refresh", auth.RefreshTokenExpiration)

	http.SetCookie(c.Writer, accessCookie)
	http.SetCookie(c.Writer, refreshCookie)

	c.JSON(http.StatusOK,
		responses.APIResponse{Message: responses.GetResponse(localizer, responses.LoginSuccess)})
}

func Logout(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	accessCookie := auth.DeleteCookie(auth.Access, "/")
	refreshCookie := auth.DeleteCookie(auth.Refresh, "/api/auth/refresh")

	http.SetCookie(c.Writer, accessCookie)
	http.SetCookie(c.Writer, refreshCookie)

	c.JSON(http.StatusOK,
		responses.APIResponse{Message: responses.GetResponse(localizer, responses.LogoutSuccess)},
	)
}

func Refresh(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	refreshSecret, err := auth.GetSecretKey(auth.Refresh)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
		return
	}

	userToken, err := auth.ValidateToken(c, auth.Refresh, refreshSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.UnauthorizedError)})
		return
	}

	accessSecret, err := auth.GetSecretKey(auth.Access)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
		return
	}

	accessToken, err := auth.GenerateToken(*userToken, accessSecret, auth.AccessTokenExpiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
		return
	}

	accessCookie := auth.CreateCookie(auth.Access, accessToken, "/", auth.AccessTokenExpiration)
	http.SetCookie(c.Writer, accessCookie)

	c.JSON(http.StatusOK, 
	responses.APIResponse{Message: responses.GetResponse(localizer, responses.TokenRenewed)},)
}
