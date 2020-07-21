package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	HostIp      string `yaml:"hostIp"`
	BlenderPath string `yaml:"blenderPath"`
}

func UnmarshalConfig() Config {
	config := Config{}
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("[ERROR] Failed to open config file\n")
	}

	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatal("[ERROR] Failed to unmarshal config object\n")
	}
	return config
}

func (c *Config) SaveConfig() {
	marshal, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatal("[ERROR] Failed to marshal config object\n")
	}

	if err := ioutil.WriteFile("config.yaml", marshal, 0644); err != nil {
		log.Fatal("[ERROR] Failed to save config file\n")
	}
}