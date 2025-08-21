package configs

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	cfg    *config
	logger *Logger
	db     *gorm.DB
)

type config struct {
	API    APIConfig
	DB     DBConfig
	CORS   CORSConfig
	AI     AIConfig
	Stripe StripeConfig
	CSRF   struct {
		Secret string
	}
	Email      EmailConfig
	WahaConfig WahaApi
}

type APIConfig struct {
	Url       string
	Port      string
	JwtSecret string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

type AIConfig struct {
	ApiKey   string
	Endpoint string
	Model    string
}

type StripeConfig struct {
	SecretKey      string
	PublishableKey string
	SuccessURL     string
	CancelURL      string
}

type EmailConfig struct {
	ApiKey string
}

type WahaApi struct {
	Url               string
	ApiKey            string
	WebhookForwardUrl string
}

func init() {
	viper.SetDefault("api.port", "9000")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
}

func Load() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	cfg = new(config)
	cfg.API = APIConfig{
		Url:       viper.GetString("api.url"),
		Port:      viper.GetString("api.port"),
		JwtSecret: viper.GetString("api.jwt_secret"),
	}
	cfg.DB = DBConfig{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.pass"),
		Database: viper.GetString("database.name"),
	}
	cfg.CORS = CORSConfig{
		AllowOrigins:     viper.GetStringSlice("cors.allow_origins"),
		AllowMethods:     viper.GetStringSlice("cors.allow_methods"),
		AllowHeaders:     viper.GetStringSlice("cors.allow_headers"),
		ExposeHeaders:    viper.GetStringSlice("cors.expose_headers"),
		AllowCredentials: viper.GetBool("cors.allow_credentials"),
		MaxAge:           viper.GetInt("cors.max_age"),
	}
	cfg.AI = AIConfig{
		ApiKey:   viper.GetString("ai.api_key"),
		Endpoint: viper.GetString("ai.endpoint"),
		Model:    viper.GetString("ai.model"),
	}
	cfg.Stripe = StripeConfig{
		SecretKey:      viper.GetString("stripe.secret_key"),
		PublishableKey: viper.GetString("stripe.publishable_key"),
		SuccessURL:     viper.GetString("stripe.success_url"),
		CancelURL:      viper.GetString("stripe.cancel_url"),
	}
	cfg.CSRF.Secret = viper.GetString("csrf.secret")
	cfg.Email.ApiKey = viper.GetString("resend.api_key")
	cfg.WahaConfig = WahaApi{
		Url:               viper.GetString("waha.waha_api_url"),
		ApiKey:            viper.GetString("waha.waha_api_key"),
		WebhookForwardUrl: viper.GetString("waha.webhook_forward_url"),
	}
	return nil
}

func GetDB() *gorm.DB {
	return db
}

func GetCSRFSecret() string {
	return cfg.CSRF.Secret
}

func GetServerUrl() string {
	return cfg.API.Url
}

func GetServerPort() string {
	return cfg.API.Port
}

func GetJwtSecret() []byte {
	secret := cfg.API.JwtSecret
	if secret == "" {
		secret = "olha_a_foto_super_secret_key_2024" // Fallback para desenvolvimento
	}
	return []byte(secret)
}

func GetLogger(p string) *Logger {
	logger = NewLogger(p)
	return logger
}

func GetCORSConfig() CORSConfig {
	return cfg.CORS
}

func GetAIConfig() AIConfig {
	return cfg.AI
}

func GetStripeConfig() StripeConfig {
	return cfg.Stripe
}

func GetEmailApiKey() string {
	return cfg.Email.ApiKey
}

func GetWahaConfig() WahaApi {
	return cfg.WahaConfig
}

func InitDB() error {
	var err error
	db, err = InitializePostgres(cfg.DB)
	if err != nil {
		return err
	}
	return nil
}
