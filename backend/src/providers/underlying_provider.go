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
	underlyingMu              sync.RWMutex
)

// GetUnderlyingService returns the singleton UnderlyingService instance.
func GetUnderlyingService() services.UnderlyingService {
	underlyingMu.RLock()
	if underlyingServiceInstance != nil {
		defer underlyingMu.RUnlock()
		return underlyingServiceInstance
	}
	underlyingMu.RUnlock()

	underlyingOnce.Do(func() {
		underlyingMu.Lock()
		defer underlyingMu.Unlock()
		if underlyingServiceInstance == nil {
			db := database.GetDB()
			repo := repositories.NewUnderlyingRepository(db)
			underlyingServiceInstance = services.NewUnderlyingService(repo)
		}
	})

	underlyingMu.RLock()
	defer underlyingMu.RUnlock()
	return underlyingServiceInstance
}

// SetUnderlyingService allows overriding the UnderlyingService instance (e.g. for testing).
func SetUnderlyingService(svc services.UnderlyingService) {
	underlyingMu.Lock()
	defer underlyingMu.Unlock()
	underlyingServiceInstance = svc
}
