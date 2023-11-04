package role

// 特定のロールを持っているか検証します
func HasAllowedRoleID(allowedRoles []string, hasRoles []string) bool {
	for _, allowedRole := range allowedRoles {
		for _, hasRole := range hasRoles {
			if allowedRole == hasRole {
				return true
			}
		}
	}

	return false
}
