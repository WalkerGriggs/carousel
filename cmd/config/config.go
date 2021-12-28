package config

import "github.com/walkergriggs/carousel/carousel"

type Config struct {
	Server *ServerConfig
	Users  []*UserConfig
}

type ServerConfig struct {
	URI             string
	Verbose         bool
	SSL             bool
	CertificatePath string
}

type UserConfig struct {
	Username string
	Password string
	Networks []*NetworkConfig
}

type NetworkConfig struct {
	Name     string
	URI      string
	Ident    *IdentityConfig
	Channels []string
}

type IdentityConfig struct {
	Username string
	Nickname string
	Realname string
	Password string
}

func EmptyConfig() *Config {
	return &Config{
		Server: &ServerConfig{},
		Users: []*UserConfig{
			{
				Networks: []*NetworkConfig{
					{
						Ident:    &IdentityConfig{},
						Channels: []string{},
					},
				},
			},
		},
	}
}

func ConvertServerConfig(raw *Config) (*carousel.ServerConfig, error) {
	users := make([]*carousel.User, len(raw.Users))
	for i, rawUser := range raw.Users {
		userConfig, err := ConvertUserConfig(rawUser)
		if err != nil {
			return nil, err
		}

		user, err := carousel.NewUser(userConfig)
		if err != nil {
			return nil, err
		}

		users[i] = user
	}

	return &carousel.ServerConfig{
		Users:           users,
		SSLEnabled:      raw.Server.SSL,
		URI:             raw.Server.URI,
		CertificatePath: raw.Server.CertificatePath,
	}, nil
}

func ConvertUserConfig(raw *UserConfig) (*carousel.UserConfig, error) {
	networks := make([]*carousel.Network, len(raw.Networks))

	for i, rawNetwork := range raw.Networks {
		config, err := ConvertNetworkConfig(rawNetwork)
		if err != nil {
			return nil, err
		}

		networks[i], err = carousel.NewNetwork(config)
		if err != nil {
			return nil, err
		}
	}

	return &carousel.UserConfig{
		Username: raw.Username,
		Password: raw.Password,
		Networks: networks,
	}, nil
}

func ConvertNetworkConfig(raw *NetworkConfig) (*carousel.NetworkConfig, error) {
	ident := carousel.Identity{
		Username: raw.Ident.Username,
		Nickname: raw.Ident.Nickname,
		Realname: raw.Ident.Realname,
		Password: raw.Ident.Password,
	}

	return &carousel.NetworkConfig{
		Name:  raw.Name,
		URI:   raw.URI,
		Ident: &ident,
	}, nil
}
