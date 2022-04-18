package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
    "path/filepath"
	"exmo-trading/app/trader"
	"exmo-trading/app/dataserver"
)

type Config interface {
	Nothing()
}

func Load(c Config, configFileName string) error{
	filename, _ := filepath.Abs(configFileName)
    yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err := yaml.Unmarshal([]byte(data), c)
	if err != nil {
		return err
	}
	return nil
}