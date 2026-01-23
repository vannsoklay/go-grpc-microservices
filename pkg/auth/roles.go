package auth

var RolePermissions = map[string][]string{
	"USER": {
		"payment:create",
		"payment:view",
	},
	"ADMIN": {
		"payment:create",
		"payment:view",
		"payment:refund",
	},
}
