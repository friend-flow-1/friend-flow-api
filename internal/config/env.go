package config

import (
	_ "log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
)

type Config struct {
	// ScyllaDB
	ScyllaDBHost     string
	ScyllaDBKeyspace string

	// PostgreSQL
	PostgreHost     string
	PostgrePort     string
	PostgreUser     string
	PostgrePassword string
	PostgreDBName   string

	// App
	JWTSecret          string
	AccessTokenSecret  string
	RefreshTokenSecret string
	ApiPort            string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool

	// Snowflake ID generator
	SnowflakeNodeID string
}

var AppConfig *Config

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		// ScyllaDB config
		ScyllaDBHost:     GetEnv("SCYLLA_DB_HOST", "127.0.0.1:9042"),
		ScyllaDBKeyspace: GetEnv("SCYLLA_DB_KEYSPACE", "friendflow"),

		// PostgreSQL config
		PostgreHost:     GetEnv("POSTGRE_DB_HOST", "127.0.0.1"),
		PostgrePort:     GetEnv("POSTGRE_DB_PORT", "5432"),
		PostgreUser:     GetEnv("POSTGRE_DB_USER", "haxxu"),
		PostgrePassword: GetEnv("POSTGRE_DB_PASSWORD", "User123"),
		PostgreDBName:   GetEnv("POSTGRE_DB_NAME", "friendflow"),

		// Application config
		JWTSecret:          GetEnv("JWT_SECRET", "friendflowsecret"),
		AccessTokenSecret:  GetEnv("ACCESS_TOKEN_SECRET", "friendflowsecretaccess"),
		RefreshTokenSecret: GetEnv("REFRESH_TOKEN_SECRET", "friendflowsecretrefresh"),
		ApiPort:            GetEnv("API_PORT", "8080"),

		MinioEndpoint:  GetEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey: GetEnv("MINIO_ACCESS_KEY", "haxxu"),
		MinioSecretKey: GetEnv("MINIO_SECRET_KEY", "User1234"),
		MinioBucket:    GetEnv("MINIO_BUCKET", "friendflowmediauploads"),
		MinioUseSSL:    GetEnv("MINIO_USE_SSL", "false") == "true",

		SnowflakeNodeID: GetEnv("SNOWFLAKE_NODE_ID", "1"),
	}
	AppConfig = cfg
	return cfg
}

func GetEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
