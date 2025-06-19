package helpers

import (
	"context"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/constants"
)

func GetLoggedInUserKey(username string) string {
	return constants.LoggedInUserKeyPrefix + username
}

func IsLoggedIn(cacheService shared.CacheService, username string) bool {
	result, err := cacheService.Exists(context.Background(), GetLoggedInUserKey(username)).Result()
	if err != nil {
		return false
	}

	return result > 0
}

func AddLoggedInUser(cacheService shared.CacheService, username string, id uint32) {
	cacheService.Set(context.Background(), GetLoggedInUserKey(username), id, 0)
}

func RemoveLoggedInUser(cacheService shared.CacheService, username string) {
	cacheService.Del(context.Background(), GetLoggedInUserKey(username))
}
