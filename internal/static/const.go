
package static

// Status constants
var (
	CREATE     = "create"
	SUCCESS    = "success"
	FAIL       = "fail"
	PROCESSING = "processing"
	CANCEL     = "cancel"
	PENDING    = "pending"
	COMPLETED  = "completed"
	ACTIVE     = "active"
	INACTIVE   = "inactive"
)

// Layer constants
var (
	HANDLER    = "handler"
	SERVICE    = "service"
	REPOSITORY = "repository"
	MIDDLEWARE = "middleware"
)

// HTTP Methods
var (
	METHOD_GET    = "GET"
	METHOD_POST   = "POST"
	METHOD_PUT    = "PUT"
	METHOD_PATCH  = "PATCH"
	METHOD_DELETE = "DELETE"
)

// Environment constants
var (
	ENV_LOCAL       = "local"
	ENV_DEVELOPMENT = "development"
	ENV_STAGING     = "staging"
	ENV_PRODUCTION  = "production"
)

// Header constants
var (
	HEADER_CONTENT_TYPE   = "Content-Type"
	HEADER_AUTHORIZATION  = "Authorization"
	HEADER_REQUEST_ID     = "X-Request-ID"
	HEADER_CORRELATION_ID = "X-Correlation-ID"
)

// Content type constants
var (
	CONTENT_TYPE_JSON = "application/json"
	CONTENT_TYPE_XML  = "application/xml"
	CONTENT_TYPE_FORM = "application/x-www-form-urlencoded"
)

// Time format constants
var (
	DATE_FORMAT      = "2006-01-02"
	DATETIME_FORMAT  = "2006-01-02 15:04:05"
	TIMESTAMP_FORMAT = "2006-01-02T15:04:05Z07:00"
)

// Pagination constants
var (
	DEFAULT_PAGE_SIZE = 20
	MAX_PAGE_SIZE     = 100
	DEFAULT_PAGE      = 1
)

// Cache duration constants (in seconds)
var (
	CACHE_SHORT  = 300   // 5 minutes
	CACHE_MEDIUM = 1800  // 30 minutes
	CACHE_LONG   = 3600  // 1 hour
	CACHE_DAY    = 86400 // 24 hours
)

