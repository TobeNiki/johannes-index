package database

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
	DBMS     string `yaml:"DBMS"`
	USER     string `yaml:"USER"`
	PASS     string `yaml:"PASS"`
	PROTOCOL string `yaml:"PROTOCOL"`
	DBNAME   string `yaml:"DBNAME"`
	PARAME   string `yaml:"PARAME"`
}

const yamlFildePath = "./backend/database/config.yml"

// const yamlFildePath = "config.yaml" test用
func LoadConfig() DatabaseConfig {
	file, err := ioutil.ReadFile(yamlFildePath)
	if err != nil {
		log.Fatal(err)
	}
	var config = DatabaseConfig{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
