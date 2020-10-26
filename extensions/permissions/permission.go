package permissions

//IsPermitted checks member have single permission
func IsPermitted(permissionsInteger, permission int) bool {
	return permissionsInteger&permission == permission || permissionsInteger&ADMINISTRATOR == ADMINISTRATOR
}

//IsPermittedAll checks member have multiple permission
func IsPermittedAll(permissionInteger int, permissions ...int) bool {
	if permissionInteger&ADMINISTRATOR == ADMINISTRATOR {
		return true
	}
	for _, permission := range permissions {
		if permissionInteger&permission != permission {
			return false
		}
	}
	return true
}
