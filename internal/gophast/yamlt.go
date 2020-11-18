package gophast

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type YamlT struct {
	Name    string `yaml:"name"`
	Port    int32  `yaml:"port"`
	Release bool   `yaml:"release"`

	PublicDir  string `yaml:"public"`
	AssetRoute string `yaml:"asset_route"`

	Subsite string `yaml:"subsite"`

	APIRoute string `yaml:"api_route"`
	APIHost  string `yaml:"api_host"`
}

func (c *YamlT) GetConfigData() *YamlT {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}