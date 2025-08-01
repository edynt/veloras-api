package initialize

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/edynnt/veloras-api/pkg/config"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func MustLoadConfig() config.Config {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("%s: %v", msg.CannotLoadConfig, err)
	}
	return cfg
}

func LoadConfig() (cfg config.Config, err error) {
	_ = godotenv.Load()

	pathConfig := "pkg/environment/config.yaml"

	raw, err := os.ReadFile(pathConfig)
	if err != nil {
		return cfg, fmt.Errorf("%s: %w", msg.CannotReadConfig, err)
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadConfig(strings.NewReader(os.ExpandEnv(string(raw)))); err != nil {
		return cfg, fmt.Errorf("%s: %w", msg.CannotReadConfig, err)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("%s: %w", msg.CannotUnmarshalConfig, err)
	}

	return
}
