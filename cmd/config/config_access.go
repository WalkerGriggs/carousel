package config

import (
	"encoding/json"
	"io/ioutil"
)

type (
	ConfigAccess interface {
		GetStartingConfig() (*Config, error)
		GetDefaultFilename() string
	}

	PathOptions struct {
		GlobalPath string
	}

	// Config struct {
	//	Users           []*user.User
	//	CertificatePath string
	// }
)

func (o *PathOptions) GetStartingConfig() (*Config, error) {
	file, err := ioutil.ReadFile(o.GlobalPath)
	if err != nil {
		return nil, err
	}

	config := Config{}

	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (o *PathOptions) GetDefaultFilename() string {
	return o.GlobalPath
}

func ModifyFile(configAccess ConfigAccess, config Config) error {
	return WriteToFile(config, configAccess.GetDefaultFilename())
}
