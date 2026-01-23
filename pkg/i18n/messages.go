package i18n

var Messages = map[string]map[string]string{
	"UNAUTHORIZED": {
		"en": "Authentication is required to access this resource.",
		"km": "ត្រូវការការផ្ទៀងផ្ទាត់មុនពេលចូលប្រើប្រាស់។",
	},
	"FORBIDDEN": {
		"en": "You do not have permission to access this resource.",
		"km": "អ្នកមិនមានសិទ្ធិចូលប្រើប្រាស់ធនធាននេះទេ។",
	},
	"NOT_FOUND": {
		"en": "The requested resource was not found.",
		"km": "រកមិនឃើញធនធានដែលបានស្នើសុំ។",
	},
	"INTERNAL_SERVER_ERROR": {
		"en": "An unexpected error occurred. Please try again later.",
		"km": "មានកំហុសមួយកើតឡើង។ សូមព្យាយាមម្តងទៀត។",
	},
	"SERVICE_UNAVAILABLE": {
		"en": "The service is temporarily unavailable. Please try again later.",
		"km": "សេវាកម្មមិនអាចប្រើប្រាស់បានបណ្តោះអាសន្ន។",
	},
}
