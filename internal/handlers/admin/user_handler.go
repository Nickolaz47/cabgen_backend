package admin

import (
	"errors"
	"net/http"
	"slices"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllUsers(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var users []models.User
	if err := db.DB.Preload("Country").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: users})
}

func GetUserByUsername(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	username := c.Param("username")

	var user models.User
	err := db.DB.Preload("Country").Where("username = ?", username).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	if err != nil {
		c.JSON(http.StatusNotFound,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.UserNotFoundError)})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: user})
}

func CreateUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			responses.APIResponse{Error: responses.GetResponse(localizer,
				responses.UnauthorizedError)})
		return
	}

	var newUser models.AdminRegisterInput
	if errMsg, valid := validations.Validate(c, localizer, &newUser); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if !slices.Contains(models.UserRoles, newUser.UserRole) {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.InvalidUserRoleError)},
		)
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

	country, valid := validations.ValidateCountryCode(newUser.CountryCode)
	if !valid {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.CountryNotFoundError)})
		return
	}

	activatedOn := time.Now()
	user := models.User{
		Name:        newUser.Name,
		Username:    newUser.Username,
		Email:       newUser.Email,
		Password:    hashedPassword,
		CountryCode: newUser.CountryCode,
		Country:     *country,
		IsActive:    true,
		UserRole:    newUser.UserRole,
		Interest:    newUser.Interest,
		Role:        newUser.Role,
		Institution: newUser.Institution,
		CreatedBy:   userToken.Username,
		ActivatedBy: &userToken.Username,
		ActivatedOn: &activatedOn,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.RegisterCreateUserError)},
		)
		return
	}

	c.JSON(http.StatusCreated,
		responses.APIResponse{
			Data:    user.ToResponse(c),
			Message: responses.GetResponse(localizer, responses.AdminRegisterSuccess),
		},
	)
}

func UpdateUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	username := c.Param("username")

	var userToUpdate models.AdminUpdateInput
	if errMsg, valid := validations.Validate(c, localizer, &userToUpdate); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	var user models.User
	err := db.DB.Preload("Country").Where("username = ?", username).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	if err != nil {
		c.JSON(http.StatusNotFound,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.UserNotFoundError)})
		return
	}

	if userToUpdate.UserRole != nil {
		if !slices.Contains(models.UserRoles, *userToUpdate.UserRole) {
			c.JSON(http.StatusBadRequest,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.InvalidUserRoleError)},
			)
			return
		}
		user.UserRole = *userToUpdate.UserRole
	}

	if userToUpdate.Email != nil || userToUpdate.Password != nil {
		var existingUser models.User
		err := db.DB.Where("email = ? OR username = ?", userToUpdate.Email, userToUpdate.Username).
			First(&existingUser).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
			)
			return
		}

		if err == nil {
			var errMsg string
			if existingUser.Email == *userToUpdate.Email {
				errMsg = responses.RegisterEmailAlreadyExistsError
			} else if existingUser.Username == *userToUpdate.Username {
				errMsg = responses.RegisterUsernameAlreadyExistsError
			}

			c.JSON(http.StatusConflict,
				responses.APIResponse{Error: responses.GetResponse(localizer, errMsg)},
			)
			return
		}
	}

	if userToUpdate.CountryCode != nil {
		country, valid := validations.ValidateCountryCode(*userToUpdate.CountryCode)
		if !valid {
			c.JSON(http.StatusBadRequest, responses.APIResponse{
				Error: responses.GetResponse(localizer, responses.CountryNotFoundError),
			})
			return
		}
		user.Country = *country
	}

	if userToUpdate.Password != nil {
		hashedPassword, err := security.Hash(*userToUpdate.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
			)
			return
		}
		user.Password = hashedPassword
	}

	validations.ApplyAdminUpdateToUser(&user, &userToUpdate)

	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.UpdateUserError),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: user.ToResponse(c),
	})
}

func DeleteUser(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	username := c.Param("username")

	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	if err != nil {
		c.JSON(http.StatusNotFound,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.UserNotFoundError)})
		return
	}

	if err := db.DB.Where("username = ?", username).Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.APIResponse{Error: responses.GetResponse(localizer, responses.GenericInternalServerError)},
		)
		return
	}

	c.JSON(http.StatusOK,
		responses.APIResponse{Message: responses.GetResponse(localizer, responses.UserDeleted)},
	)
}
