package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Log struct {
		Level  string
		Format string
		Output string
		Core   string
	}
	Gim struct {
		Env   string
		Debug bool
	}
	Server struct {
		IP   string
		Name string
		Host string
		Port string
	}

	Connection struct {
		HTTP struct {
			TimeOut int
		}
		Postgresql struct {
			Master struct {
				DB       string
				Host     string
				User     string
				Password string
				MaxOpen  int
				MaxIdle  int
			}
			Slave struct {
				DB       string
				Host     string
				User     string
				Password string
				MaxOpen  int
				MaxIdle  int
			}
			FixedReadInstance string
		}
	}
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config.Newconfig: load file config failed: %w", err)
	}

	var conf Config

	if err := v.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("config.Newconfig: unmarshal config failed: %w", err)
	}

	return &conf, nil
}
