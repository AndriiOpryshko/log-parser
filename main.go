package main

import (
	"log-parser/db"
	"log-parser/logserv"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)
	// Parse yml config from ./config.yml
	InitConfig()


	// init db service
	conStr := GetMongoConnectionString()
	dbAuth, dbName, login, pass := GetDbCred()
	dbService := db.NewMongoService(conStr, dbAuth, dbName, login, pass)

	// init log parser service
	logsConf := GetParseLogsConfig()
	logService := logserv.NewLogService(dbService, ConvToLogConf(logsConf))

	// Running
	logService.Run()
}
