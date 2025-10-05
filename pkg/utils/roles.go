package utils

import "github.com/ranggakrisnaa/go-fiber-starter/database/entities"

func RoleNames(roles []entities.Role) []string {
	names := make([]string, len(roles))
	for i, r := range roles {
		names[i] = r.GetName()
	}
	return names
}
