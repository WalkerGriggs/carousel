package config

import (
	// "encoding/json"
	// "io/ioutil"

	// "github.com/walkergriggs/carousel/pkg/channel"
	"github.com/walkergriggs/carousel/pkg/identity"
	"github.com/walkergriggs/carousel/pkg/network"
	"github.com/walkergriggs/carousel/pkg/server"
	"github.com/walkergriggs/carousel/pkg/user"
)

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
	Network  *NetworkConfig
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
				Network: &NetworkConfig{
					Ident:    &IdentityConfig{},
					Channels: []string{},
				},
			},
		},
	}
}

func ConvertServerConfig(raw *Config) (*server.Config, error) {
	users := make([]*user.User, len(raw.Users))
	for i, rawUser := range raw.Users {
		userConfig, err := ConvertUserConfig(rawUser)
		if err != nil {
			return nil, err
		}

		user, err := user.New(userConfig)
		if err != nil {
			return nil, err
		}

		users[i] = user
	}

	return &server.Config{
		Users:           users,
		SSLEnabled:      raw.Server.SSL,
		URI:             raw.Server.URI,
		CertificatePath: raw.Server.CertificatePath,
	}, nil
}

func ConvertUserConfig(raw *UserConfig) (*user.Config, error) {
	networkConfig, err := ConvertNetworkConfig(raw.Network)
	if err != nil {
		return nil, err
	}

	network, err := network.New(networkConfig)
	if err != nil {
		return nil, err
	}

	return &user.Config{
		Username: raw.Username,
		Password: raw.Password,
		Network:  network,
	}, nil
}

func ConvertNetworkConfig(raw *NetworkConfig) (*network.Config, error) {
	ident := identity.Identity{
		Username: raw.Ident.Username,
		Nickname: raw.Ident.Nickname,
		Realname: raw.Ident.Realname,
		Password: raw.Ident.Password,
	}

	// channels := make([]*channel.Channel, len(raw.Channels))
	// for i, name := range raw.Channels {
	//	channels[i] = &channel.Channel{
	//		Name: name,
	//	}
	// }

	return &network.Config{
		Name: raw.Name,
		URI:  raw.URI,
		// Channels: channels,
		Ident: &ident,
	}, nil
}
