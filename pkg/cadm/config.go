package cadm

import (
	"path"
	"strings"

	"github.com/spf13/viper"

	"github.com/walkergriggs/carousel/carousel"
)

type (
	Config struct {
		LoggingLevel string
		Addrs        *AdvertiseAddrs
		Server       *ServerConfig
	}

	AdvertiseAddrs struct {
		IRC  string
		HTTP string
	}

	ServerConfig struct {
		SSL             bool
		CertificatePath string
	}
)

func ConvertServerConfig(raw *Config) (*carousel.ServerConfig, error) {
	return &carousel.ServerConfig{
		SSLEnabled:      raw.Server.SSL,
		URI:             raw.Addrs.IRC,
		CertificatePath: raw.Server.CertificatePath,
	}, nil
}

// unmarshalConfig reads in the config file from either the default directory
// or the given path and returns the extracted values in a Config struct.
func unmarshalConfig(p string) (*Config, error) {
	viper.BindEnv("log_level")
	viper.AutomaticEnv()

	configPath, configFile := path.Split(p)
	configName := strings.TrimSuffix(configFile, path.Ext(p))

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("OS")

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{
		Addrs:  &AdvertiseAddrs{},
		Server: &ServerConfig{},
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
