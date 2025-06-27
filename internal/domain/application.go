package domain

// AppConfig - represents application configuration
type AppConfig struct {
	Environment    string `env:"ENVIRONMENT" validate:"required"`
	AllowedOrigins string `env:"ALLOWED_ORIGINS" validate:"required"`
	Port           string `env:"PORT" validate:"required"`
	DSN            string `env:"DSN" validate:"required"`
}
