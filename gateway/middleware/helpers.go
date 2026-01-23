package middleware

func hasPermissions(userPerms, requiredPerms []string) bool {
	permSet := make(map[string]struct{}, len(userPerms))
	for _, p := range userPerms {
		permSet[p] = struct{}{}
	}

	for _, rp := range requiredPerms {
		if _, ok := permSet[rp]; !ok {
			return false
		}
	}

	return true
}
