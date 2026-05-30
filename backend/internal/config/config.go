package config

import (
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port string
	Env  string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret        string
	Expiry        time.Duration
	RefreshExpiry time.Duration
}

type FrontendConfig struct {
	URL      string
	AdminURL string
}

type KUNConfig struct {
	// AppKey is the "App-Key" request header required by KUN.
	AppKey string
	// ApiVersion is the "Api-Version" request header required by KUN (e.g. "2").
	ApiVersion string
	// CustomerNo is the "Customer-No" request header required by KUN.
	CustomerNo string
	// PrivateKey is the RSA private key PEM used for SHA256withRSA request signing.
	PrivateKey string
	// PublicKey is the RSA public key PEM used for verifying KUN signatures (webhook/response).
	PublicKey       string
	BaseURL         string
	RegionCode      string
	Timeout         time.Duration
	WebhookTimeDiff time.Duration
	MockEnabled     bool
}

type Config struct {
	App      AppConfig
	DB       DBConfig
	JWT      JWTConfig
	Frontend FrontendConfig
	KUN      KUNConfig
}

func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "3306")
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "motewallet")
	viper.SetDefault("JWT_SECRET", "motewallet-jwt-secret-change-in-production")
	viper.SetDefault("JWT_EXPIRY", "24h")
	viper.SetDefault("JWT_REFRESH_EXPIRY", "168h")
	viper.SetDefault("FRONTEND_URL", "http://localhost:3000")
	viper.SetDefault("ADMIN_URL", "http://localhost:3001")
	viper.SetDefault("KUN_API_BASE_URL", "https://api.kun.global")
	viper.SetDefault("KUN_REGION_CODE", "KUN_PL")
	viper.SetDefault("KUN_API_TIMEOUT", "30s")
	viper.SetDefault("KUN_WEBHOOK_TIME_DIFF", "5m")
	viper.SetDefault("KUN_MOCK_ENABLED", false)
	viper.SetDefault("KUN_API_VERSION", "2")

	// Ignore error if .env file doesn't exist — env vars still work
	_ = viper.ReadInConfig()

	jwtExpiry, err := time.ParseDuration(viper.GetString("JWT_EXPIRY"))
	if err != nil {
		jwtExpiry = 24 * time.Hour
	}

	jwtRefreshExpiry, err := time.ParseDuration(viper.GetString("JWT_REFRESH_EXPIRY"))
	if err != nil {
		jwtRefreshExpiry = 168 * time.Hour
	}

	kunTimeout, err := time.ParseDuration(viper.GetString("KUN_API_TIMEOUT"))
	if err != nil {
		kunTimeout = 30 * time.Second
	}

	kunWebhookTimeDiff, err := time.ParseDuration(viper.GetString("KUN_WEBHOOK_TIME_DIFF"))
	if err != nil {
		kunWebhookTimeDiff = 5 * time.Minute
	}

	cfg := &Config{
		App: AppConfig{
			Port: viper.GetString("APP_PORT"),
			Env:  viper.GetString("APP_ENV"),
		},
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		JWT: JWTConfig{
			Secret:        viper.GetString("JWT_SECRET"),
			Expiry:        jwtExpiry,
			RefreshExpiry: jwtRefreshExpiry,
		},
		Frontend: FrontendConfig{
			URL:      viper.GetString("FRONTEND_URL"),
			AdminURL: viper.GetString("ADMIN_URL"),
		},
		KUN: KUNConfig{
			AppKey:          viper.GetString("KUN_APP_KEY"),
			ApiVersion:      viper.GetString("KUN_API_VERSION"),
			PrivateKey:      viper.GetString("KUN_PRIVATE_KEY"),
			PublicKey:        viper.GetString("KUN_PUBLIC_KEY"),
			BaseURL:         viper.GetString("KUN_API_BASE_URL"),
			RegionCode:      viper.GetString("KUN_REGION_CODE"),
			CustomerNo:      viper.GetString("KUN_CUSTOMER_NO"),
			Timeout:         kunTimeout,
			WebhookTimeDiff: kunWebhookTimeDiff,
			MockEnabled:     viper.GetBool("KUN_MOCK_ENABLED"),
		},
	}

	return cfg, nil
}

func (db *DBConfig) DSN() string {
	return db.User + ":" + db.Password + "@tcp(" + db.Host + ":" + db.Port + ")/" + db.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
}
