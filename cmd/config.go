package cmd

import (
	"github.com/spf13/viper"
	"github.com/walkergriggs/carousel/pkg/server"
)

func unmarshalConfig() (*server.Server, error) {
	var s server.Server

	err := viper.Unmarshal(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
