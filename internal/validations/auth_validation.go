package validations

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/gin-gonic/gin"
)

const UserTokenKey = "User"

func GetUserTokenFromContext(c *gin.Context) (*models.UserToken, bool) {
	rawUserToken, exists := c.Get(UserTokenKey)
	if !exists {
		return nil, false
	}

	userToken, ok := rawUserToken.(*models.UserToken)
	if !ok {
		return nil, false
	}

	return userToken, true
}
