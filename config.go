package main

import (
	"path/filepath"
	"log-parser/log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

// main configs
type Config struct {
	MongoDbConfig  MongoDBConf   `yaml:"mongodb"`
	ParseLogsConfig []ParseLogConf `yaml:"logs"`
}

// mongo config
type MongoDBConf struct {
	Ip       string `yaml:"ip"`
	Port     string `yaml:"port"`
	DbName   string `yaml:"db_name"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
}

// log config
type ParseLogConf struct {
	AbsPath string `yaml:"abs_path"`
	Type string `yaml:"log_type"`
}


// Config global object
var config *Config = &Config{}


// Init all configs from config.yml
func InitConfig() {
	absPath, _ := filepath.Abs("./config.yml")
	log.Infof("Start init config. Config path %s", absPath)

	yamlFile, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Criticalf("yamlFile.Get err %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Criticalf("Unmarshal: %v", err)
	}

	log.Debugf("%+v", *config)
	log.Info("Success init config")
}
