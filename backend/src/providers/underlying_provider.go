package providers

import (
	"sync"

	"backend/src/database"
	"backend/src/repositories"
	"backend/src/services"
)

var (
	underlyingServiceInstance services.UnderlyingService
	underlyingOnce            sync.Once
)

// GetUnderlyingService returns the singleton UnderlyingService instance.
func GetUnderlyingService() services.UnderlyingService {
	underlyingOnce.Do(func() {
		db := database.GetDB()
		repo := repositories.NewUnderlyingRepository(db)
		underlyingServiceInstance = services.NewUnderlyingService(repo)
	})
	return underlyingServiceInstance
}

// SetUnderlyingService allows overriding the UnderlyingService instance (e.g. for testing).
func SetUnderlyingService(svc services.UnderlyingService) {
	underlyingServiceInstance = svc
}
