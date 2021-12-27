package cmd

import (
	"github.com/spf13/viper"

	"github.com/walkergriggs/carousel/carousel"
)

func unmarshalConfig() (*carousel.Server, error) {
	var s carousel.Server

	err := viper.Unmarshal(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
