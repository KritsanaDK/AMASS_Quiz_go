
package static

// Success codes
var (
	CODE_SUCCESS = "0000"
)

// Client error codes (4xxx)
var (
	CODE_BAD_REQUEST         = "4000"
	CODE_UNAUTHORIZED        = "4001"
	CODE_FORBIDDEN           = "4003"
	CODE_NOT_FOUND           = "4004"
	CODE_METHOD_NOT_ALLOWED  = "4005"
	CODE_CONFLICT            = "4009"
	CODE_VALIDATION_ERROR    = "4022"
	CODE_RATE_LIMIT_EXCEEDED = "4029"
)

// Server error codes (5xxx)
var (
	CODE_INTERNAL_ERROR      = "5000"
	CODE_SERVICE_UNAVAILABLE = "5003"
	CODE_GATEWAY_TIMEOUT     = "5004"
)

// Database error codes (5xxx)
var (
	CODE_DATABASE_ERROR      = "5100"
	CODE_DATABASE_TIMEOUT    = "5101"
	CODE_DATABASE_CONSTRAINT = "5102"
)

// External service error codes (5xxx)
var (
	CODE_EXTERNAL_SERVICE_ERROR = "5200"
	CODE_AWS_SERVICE_ERROR      = "5201"
	CODE_SECRET_MANAGER_ERROR   = "5202"
)

// Error messages
var ErrorMessages = map[string]string{
	CODE_SUCCESS:                "Success",
	CODE_BAD_REQUEST:            "Bad request",
	CODE_UNAUTHORIZED:           "Unauthorized",
	CODE_FORBIDDEN:              "Forbidden",
	CODE_NOT_FOUND:              "Resource not found",
	CODE_METHOD_NOT_ALLOWED:     "Method not allowed",
	CODE_CONFLICT:               "Resource conflict",
	CODE_VALIDATION_ERROR:       "Validation error",
	CODE_RATE_LIMIT_EXCEEDED:    "Rate limit exceeded",
	CODE_INTERNAL_ERROR:         "Internal server error",
	CODE_SERVICE_UNAVAILABLE:    "Service unavailable",
	CODE_GATEWAY_TIMEOUT:        "Gateway timeout",
	CODE_DATABASE_ERROR:         "Database error",
	CODE_DATABASE_TIMEOUT:       "Database timeout",
	CODE_DATABASE_CONSTRAINT:    "Database constraint violation",
	CODE_EXTERNAL_SERVICE_ERROR: "External service error",
	CODE_AWS_SERVICE_ERROR:      "AWS service error",
	CODE_SECRET_MANAGER_ERROR:   "Secret manager error",
}

