package settings

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Conf struct {
	DbHost     string `yaml:"db_host"`
	DbName     string `yaml:"db_name"`
	DbUser     string `yaml:"db_user"`
	DbPassword string `yaml:"db_password"`
	DbPort     string `yaml:"db_port"`
}

func GetConf(filePath string) *Conf {
	var c *Conf
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}
	return c
}
