package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CABGenOrg/cabgen_backend/internal/events"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func UserRegisteredHandler(emailService services.EmailService) events.HandlerFunc {
	return func(ctx context.Context, payload []byte) error {
		var data events.UserRegisteredPayload

		if err := json.Unmarshal(payload, &data); err != nil {
			return fmt.Errorf("invalid payload: %w", err)
		}

		return emailService.SendActivationUserEmail(
			ctx, data.RegisteredUsername,
		)
	}
}
