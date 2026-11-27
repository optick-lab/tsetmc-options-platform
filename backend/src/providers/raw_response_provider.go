package providers

import (
	"sync"

	"backend/src/database"
	"backend/src/repositories"
	"backend/src/services"
)

var (
	rawResponseServiceInstance services.RawResponseService
	rawResponseOnce            sync.Once
	rawResponseMu              sync.RWMutex
)

// GetRawResponseService returns the singleton RawResponseService instance.
func GetRawResponseService() services.RawResponseService {
	rawResponseMu.RLock()
	if rawResponseServiceInstance != nil {
		defer rawResponseMu.RUnlock()
		return rawResponseServiceInstance
	}
	rawResponseMu.RUnlock()

	rawResponseOnce.Do(func() {
		rawResponseMu.Lock()
		defer rawResponseMu.Unlock()
		if rawResponseServiceInstance == nil {
			db := database.GetDB()
			repo := repositories.NewRawResponseRepository(db)
			rawResponseServiceInstance = services.NewRawResponseService(repo)
		}
	})

	rawResponseMu.RLock()
	defer rawResponseMu.RUnlock()
	return rawResponseServiceInstance
}

// SetRawResponseService allows overriding the RawResponseService instance (e.g. for testing).
func SetRawResponseService(svc services.RawResponseService) {
	rawResponseMu.Lock()
	defer rawResponseMu.Unlock()
	rawResponseServiceInstance = svc
}
