package public

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var newUser models.RegisterInput
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingByEmail, existingByUsername models.User
	if err := db.DB.Where("email = ?", newUser.Email, newUser.Username).First(&existingByEmail).Error; err == nil {
		c.JSON(http.StatusConflict,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.RegisterEmailAlreadyExistsError)},
		)
		return
	}

	if err := db.DB.Where("username = ?", newUser.Email, newUser.Username).First(&existingByUsername).Error; err == nil {
		c.JSON(http.StatusConflict,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.RegisterUsernameAlreadyExistsError)},
		)
		return
	}

	hashedPassword, err := security.Hash(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.RegisterHashPasswordError)},
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
