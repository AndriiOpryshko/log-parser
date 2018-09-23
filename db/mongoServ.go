package db

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"log-parser/log"
	"time"
)

const (
	logC = "logs"
)

type mongoService struct {
	session *mgo.Session
	dbName  string
}

func NewMongoService(conStr, dbAuth, dbName, userName, pas string) *mongoService {
	log.Info("Init mongo db")

	session, err := mgo.Dial(conStr)

	if err != nil {
		log.Criticalf("Mongo: %s", err)
	}

	if session.Ping() != nil {
		log.Criticalf("Mongo: %s", err)
	}

	if err := session.DB(dbAuth).Login(userName, pas); err != nil {
		log.Criticalf("Mongo: %s", err)
	} else {
		log.Info("Mongo: login success")
	}

	log.Info("Init mongo db complete")

	return &mongoService{
		session: session,
		dbName:  dbName,
	}
}

func (self mongoService) WrightToCollection(name, colName string, obj interface{}) error {
	session := self.session.Copy()
	defer session.Close()

	c := session.DB(self.dbName).C(colName)
	err := c.Insert(obj)
	if err != nil {
		if mgo.IsDup(err) {
			errMsg := fmt.Sprintf("Mongo: %s already exists: %+v", name, obj)
			return errors.New(errMsg)
		}
		errMsg := fmt.Sprintf("Mongo: Failed to insert %s: %v", name, err)
		return errors.New(errMsg)
	}

	log.Debugf("Mongo: %s wrights successful. Value: %+v", name, obj)
	return nil
}

type Log struct {
	Time   time.Time `json:"log_time"`   // (Date) - время записи из лог файла
	Msg    string `json:"log_msg"`    // (String) - текст сообщения из лог файла
	Path   string `json:"file_name"`  // (String) - путь к файлу из которого получено сообщение
	Format string `json:"log_format"` // (String) - формат лога (first_format | second_format)
}

func (self mongoService) WrightLog(time time.Time, msg, path, format string) error {
	log := Log{
		Time:time,
		Msg:msg,
		Path:path,
		Format:format,
	}
	error := self.WrightToCollection("log", logC, log)

	return error
}
