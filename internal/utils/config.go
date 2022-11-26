package utils

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Env                     string `mapstructure:"ENV"`
	DBDriver                string `mapstructure:"DB_DRIVER"`
	DBUrl                   string `mapstructure:"DATABASE_URL"`
	Host                    string `mapstructure:"HOST"`
	Port                    int    `mapstructure:"PORT"`
	ServerAddress           string `mapstructure:"SERVER_ADDRESS"`
	FrontendAddress         string `mapstructure:"FRONTEND_ADDRESS"`
	JwtSecretKey            string `mapstructure:"JWT_SECRET_KEY"`
	AccessTokenExpiredTime  int32  `mapstructure:"ACCESS_TOKEN_EXPIRED_TIME"`
	RefreshTokenExpiredTime int32  `mapstructure:"REFRESH_TOKEN_EXPIRED_TIME"`

	FBKey    string `mapstructure:"FB_KEY"`
	FBSecret string `mapstructure:"FB_SECRET"`
	GGKey    string `mapstructure:"GLE_KEY"`
	GGSecret string `mapstructure:"GLE_SECRET"`

	SendgridApiKey string `mapstructure:"SENDGRID_API_KEY"`
	SendgridEmail  string `mapstructure:"SENDGRID_EMAIL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	if config.Env == "PROD" {
		config.FBKey = os.Getenv("FB_KEY")
		config.FBSecret = os.Getenv("FB_SECRET")
		config.GGKey = os.Getenv("GLE_KEY")
		config.GGSecret = os.Getenv("GLE_SECRET")
		config.SendgridApiKey = os.Getenv("SENDGRID_API_KEY")
		config.SendgridEmail = os.Getenv("SENDGRID_EMAIL")
	}

	return
}
