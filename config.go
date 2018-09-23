package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log-parser/log"
	"log-parser/logserv"
	"path/filepath"
)

// main configs
type Config struct {
	MongoDbConfig   *MongoDBConf    `yaml:"mongodb"`
	ParseLogsConfig []*ParseLogConf `yaml:"logs"`
}

// mongo config
type MongoDBConf struct {
	Ip       string `yaml:"ip"`
	Port     string `yaml:"port"`
	DbAuth   string `yaml:"db_auth"`
	DbName   string `yaml:"db_name"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
}

// log config
type ParseLogConf struct {
	AbsPath string `yaml:"path"`
	Type    string `yaml:"log_type"`
}

// getter abspath
func (self *ParseLogConf) GetAbsPath() string {
	return self.AbsPath
}

// getter log type
func (self *ParseLogConf) GetLogType() string {
	return self.Type
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

	log.Debugf("Mongo conf %+v", *config.MongoDbConfig)
	for i, parseLogConf := range config.ParseLogsConfig {
		log.Debugf("Parse log %d is %+v", i+1, *parseLogConf)
	}
	log.Info("Success init config")
}

// Get mongo connection string
func GetMongoConnectionString() string {
	conStr := fmt.Sprintf("%s:%s", config.MongoDbConfig.Ip, config.MongoDbConfig.Port)
	log.Debugf("Mondo connection string: %s", conStr)
	return conStr
}

// Get db credentials
func GetDbCred() (dbauth, dbname, userName, password string) {
	dbauth = config.MongoDbConfig.DbAuth
	dbname = config.MongoDbConfig.DbName
	userName = config.MongoDbConfig.UserName
	password = config.MongoDbConfig.Password
	return
}

// Get parse logs config
func GetParseLogsConfig() []*ParseLogConf {
	return config.ParseLogsConfig
}

// Convert logs config from  config to  logs config of log parsing service
func ConvToLogConf(parsedLogConfig []*ParseLogConf) []logserv.LogConf {
	logConf := make([]logserv.LogConf, len(parsedLogConfig))

	for i, v := range parsedLogConfig {
		// if path was relative
		v.AbsPath, _ = filepath.Abs(v.AbsPath)

		logConf[i] = logserv.LogConf(v)
	}

	return logConf
}
