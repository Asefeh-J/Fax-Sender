package utilities

const (
	APP_NAME              string = "Print2Fax"
	APP_EXEC_FILE_NAME    string = "fax_sender"
	APP_EXEC_UI_FILE_NAME string = "fax_sender_ui"
	CONFIG_FILE_NAME      string = "config.yaml"
	DEFAULT_LISTEN_PORT   int    = 11111
	CHARS                 string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	API_PATHS             string = "/api/v1"
	SECRET_KEY            string = "FAX_SENDER"
	SETTINGS_FILE_NAME    string = "settings.bin"
	ENCRYPTION_KEY        string = "0123456789012345" // the key size should be at least 16
	SEPARATOR             string = "======================================"
	JSON_CONTENT_TYPE     string = "application/json"
	LOCALHOST             string = "127.0.0.1"
	HTTP_SCHEMA           string = "http"

	WITH_COVER    string = "1"
	WITHOUT_COVER string = "0"

	WITH_PRINT    string = "1"
	WITHOUT_PRINT string = "0"

	TWO_GB_SIZE int64 = 2_147_483_648

	ERROR_CODE_NOT_ENOUGH_ARGUMENT                int = -1
	ERROR_CODE_UNKOWN_COMMAND                     int = -2
	ERROR_CODE_DEPLOY_LINUX_VERSION_NOT_FOUNDED   int = -3
	ERROR_CODE_DEPLOY_WINDOWS_VERSION_NOT_FOUNDED int = -4
	ERROR_CODE_INIT_LOG_ERROR                     int = -5
	ERROR_CDOE_INVALID_FILE_EXTENSION             int = -6
	ERROR_CODE_WORKING_DIR_NOT_FOUND              int = -7
	ERROR_CODE_IN_INIT_FILE                       int = -8
	ERROR_CODE_VERSION_FILE_NOT_FOUND             int = -9
)
