package configs

import (
	"github.com/jinzhu/configor"
)

// Database configuration
type ConfigYaml struct {
	Mysql struct {
		UserName string `default:"root"`
		PassWord string `required:"true"`
		Host     string `required:"true"`
		Port     string `default:"3306"`
		DataBase string `required:"true"`
		CharSet  string `default:"utf8"`
	}
}

var Config = &ConfigYaml{}

func LoadConfig(file string) {
	configor.Load(Config, file)
}
