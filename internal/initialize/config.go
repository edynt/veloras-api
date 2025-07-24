package initialize

import (
	"fmt"
	"os"
	"strings"

	"github.com/edynnt/veloras-api/pkg/config"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadConfig() (cfg config.Config, err error) {
	_ = godotenv.Load()

	// Read raw YAML file
	raw, err := os.ReadFile("pkg/environment/local.yaml")
	if err != nil {
		return cfg, fmt.Errorf("cannot read yaml: %w", err)
	}

	// Replace ${ENV_VAR} with actual env
	expanded := os.ExpandEnv(string(raw))

	v := viper.New()
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config from expanded string
	err = v.ReadConfig(strings.NewReader(expanded))
	if err != nil {
		return cfg, fmt.Errorf("failed to read config: %w", err)
	}

	// Unmarshal into your struct
	err = v.Unmarshal(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}
