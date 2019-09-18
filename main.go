package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.carousel")

	var carousel Server

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&carousel); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on ", carousel.URI.Format())

	carousel.Serve()
}
