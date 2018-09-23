package main

import (
	"log-parser/db"
	"time"
)

func main() {
	InitConfig()
	conStr := GetMongoConnectionString()

	dbAuth, dbName, login, pass := GetDbCred()

	dbService := db.NewMongoService(conStr, dbAuth, dbName, login, pass)


	dbService.WrightLog(time.Now(), "testmsg", "testpath", "testformat")
}
