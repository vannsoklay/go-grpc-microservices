package domain

var RolePermissions = map[string][]string{
	"ADMIN": {
		"PermProductCreate",
		"PermProductRead",
		"PermProductUpdate",
		"PermProductDelete",

		"PermPaymentCreate",
		"PermPaymentRead",
		"PermPaymentRefund",

		"PermOrderCreate",
		"PermOrderRead",
	},

	"MERCHANT": {
		"PermProductCreate",
		"PermProductRead",
		"PermProductUpdate",

		"PermPaymentRead",
		"PermOrderRead",
	},

	"USER": {
		"PermProductRead",
		"PermPaymentCreate",
		"PermPaymentRead",
		"PermOrderCreate",
		"PermOrderRead",
	},
}
