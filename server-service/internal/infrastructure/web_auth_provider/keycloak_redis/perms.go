package keycloak_redis

func (p *Provider) IsUserHaveRoles(roles []string, userRoles []string) bool {
	if len(roles) == 0 {
		return true
	}

	for _, role := range roles {
		if IsArrayContains(userRoles, role) {
			return true
		}
	}
	return false
}

func IsArrayContains(arr []string, str string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}
