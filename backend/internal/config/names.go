package config

// Names config const.
const (
	File        = "CONGIG_FILE"  // Directory and name config file.
	LogLevel    = "LOG_LEVEL"    // Application logging level (debug/info/warn/error)
	DBHost      = "DB_HOST"      // Database host address
	DBPort      = "DB_PORT"      // Database port
	DBUser      = "DB_USER"      // Database username
	DBPassword  = "DB_PASSWORD"  // Database password
	DBName      = "DB_NAME"      // Database name
	DBSSL       = "DB_SSLMODE"   // Database SSL mode (disable/require/verify)
	HTTPAddress = "HTTP_ADDRESS" // Address for http.
	HTTPPort    = "HTTP_PORT"    // Port for http.
)

// Log level const.
const (
	LogLevelDebug = "DEBUG" // DEBUG
	LogLevelInfo  = "INFO"  // INFO
	LogLevelWarn  = "WARN"  // WARN
	LogLevelError = "ERROR" // ERROR
)
