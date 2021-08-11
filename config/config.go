package config

import (
	"encoding/json"
	"io/ioutil"
)

func LoadConfigFromFile(file string) (cfg *Config, err error) {
	result := &Config{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		return result, err
	}

	return result, nil
}
