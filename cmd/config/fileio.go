package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteToFile(config Config, filename string) error {
	content, err := Encode(config)
	if err != nil {
		return err
	}

	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(filename, content, 0600); err != nil {
		return err
	}

	return nil
}

func Encode(config Config) ([]byte, error) {
	return json.MarshalIndent(config, "", "\t")
}
