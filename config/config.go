package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Tunnel struct {
	Remote string `yaml:"remote"`
	Local  string `yaml:"local"`
}

type Cfg struct {
	Uri     string   `yaml:"uri"`
	Pkey    string   `yaml:"pkey,omitempty"`
	Tunnels []Tunnel `yaml:"tunnels,omitempty"`
}

var Config []Cfg

func Init(filename string) (err error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		err = fmt.Errorf("读取配置文件%s错误，%v", filename, err)
		return
	}

	err = yaml.Unmarshal(yamlFile, &Config)
	// err = yaml.Unmarshal(yamlFile, &resultMap)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return
}
