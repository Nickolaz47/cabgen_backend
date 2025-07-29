package public

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var newUser models.RegisterInput
	if errMsg, valid := validations.ValidateRegisterInput(c, localizer, &newUser); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	var count int64
	if err := db.DB.Model(&models.User{}).
		Where("email = ? OR username = ?", newUser.Email, newUser.Username).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	if count > 0 {
		var existingUser models.User
		if err := db.DB.Where(`email = ? OR username = ?`, newUser.Email, newUser.Username).
			First(&existingUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
			)
			return
		}
		
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

func Login(c *gin.Context) {}

func Logout(c *gin.Context) {}

func Refresh(c *gin.Context) {}
