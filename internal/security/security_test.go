package security_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	hasher := security.NewPasswordHasher()

	t.Run("Success", func(t *testing.T) {
		result, err := hasher.Hash("password")

		assert.NotEmpty(t, result)
		assert.Len(t, result, 60)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		_, err := hasher.Hash("passwordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpassword")

		assert.Error(t, err)
	})
}

func TestCheckPassword(t *testing.T) {
	hasher := security.NewPasswordHasher()

	t.Run("Success", func(t *testing.T) {
		hashPassword := "$2a$10$uweyzIqC8Zu3jJsaH82iMOivIRFgFqCKjL.3SsGcR6ykTq2nsHUAW"
		password := "password"

		err := hasher.CheckPassword(hashPassword, password)

		assert.NoError(t, err, "expected passwords to be equal")
	})

	t.Run("Error", func(t *testing.T) {
		hashPassword := "$2a$10$uweyzIqC8Zu3jJsaH82iMOivIRFgFqCKjL.3SsGcR6ykTq2nsHUA"
		password := "password"

		err := hasher.CheckPassword(hashPassword, password)

		assert.Error(t, err, "expected passwords to be different")
	})
}
