package public

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var newUser models.RegisterInput
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.User
	if err := db.DB.Where("email = ? OR username = ?", newUser.Email, newUser.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email or username already in use."})
		return
	}

	hashedPassword, err := security.Hash(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Some error occurrs. Try again later."})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user. Try again later."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user.ToResponse(),
		"message": "User created successfully. Wait for an administrator to activate it."})
}

func Login(c *gin.Context) {}

func Logout(c *gin.Context) {}

func Refresh(c *gin.Context) {}
