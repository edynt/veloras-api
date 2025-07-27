package initialize

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/edynnt/veloras-api/pkg/config"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func MustLoadConfig() config.Config {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("fatal: cannot load config: %v", err)
	}
	return cfg
}

func LoadConfig() (cfg config.Config, err error) {
	_ = godotenv.Load()

	raw, err := os.ReadFile("pkg/environment/config.yaml")
	if err != nil {
		return cfg, fmt.Errorf("cannot read config file: %w", err)
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadConfig(strings.NewReader(os.ExpandEnv(string(raw)))); err != nil {
		return cfg, fmt.Errorf("failed to read config: %w", err)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return
}
