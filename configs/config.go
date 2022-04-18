package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
)

type Config interface {
	Nothing()
}

func Load(c Config, fileName string) error{
    yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}
	fmt.Println(c)
	return nil
}