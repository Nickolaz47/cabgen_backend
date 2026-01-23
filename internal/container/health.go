package container

import "github.com/CABGenOrg/cabgen_backend/internal/handlers/public/health"

func BuildHealthHandler() *health.HealthHandler {
	return health.NewHealthHandler()
}
