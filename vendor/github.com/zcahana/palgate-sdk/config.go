package palgate

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	// Config keys
	configKeyServerAddress = "serverAddress"
	configKeyGateID        = "gateID"
	configKeyAuthToken     = "authToken"

	// Environment variables
	configEnvServerAddress = "PALGATE_SERVER_ADDRESS"
	configEnvGateID        = "PALGATE_GATE_ID"
	configEnvAuthToken     = "PALGATE_AUTH_TOKEN"

	// Default values
	configDefaultServerAddress = "api1.pal-es.com"
)

type Config struct {
	ServerAddress string
	GateID        string
	AuthToken     string
}

func InitConfig() (*Config, error) {
	viper.SetDefault(configKeyServerAddress, configDefaultServerAddress)

	viper.BindEnv(configKeyServerAddress, configEnvServerAddress)
	viper.BindEnv(configKeyGateID, configEnvGateID)
	viper.BindEnv(configKeyAuthToken, configEnvAuthToken)

	viper.SetConfigName(".palgate")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, isNotFound := err.(viper.ConfigFileNotFoundError); !isNotFound {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	}

	return &Config{
		ServerAddress: viper.GetString(configKeyServerAddress),
		GateID:        viper.GetString(configKeyGateID),
		AuthToken:     viper.GetString(configKeyAuthToken),
	}, nil
}

func (conf *Config) Validate() error {
	if conf.ServerAddress == "" {
		return fmt.Errorf("missing server address")
	}
	if conf.GateID == "" {
		return fmt.Errorf("missing gate ID")
	}
	if conf.AuthToken == "" {
		return fmt.Errorf("missing auth token")
	}

	return nil
}
