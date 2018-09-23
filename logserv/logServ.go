package logserv

import (
	"log-parser/log"
	"time"
	"io"
	"os"
	"bufio"
	"strings"
	"github.com/go-errors/errors"
	"fmt"
)

const (
	firstFormat = "first_format"
	secondFormat = "second_format"
)

var dateFormats = map[string]string {
	firstFormat:"Jan 2, 2006 at 3:04:05pm (MST)",
	secondFormat: time.RFC3339,
}

type DbService interface {
	WrightLog(time time.Time, msg, path, format string) error
}

type LogConf interface {
	GetAbsPath() string
	GetLogType() string
}

type LogLine struct {
	Time   time.Time
	Msg    string
	Path   string
	Format string
}

type logParserService struct {
	dbService DbService
	logsConf  []LogConf
	LogLines  chan *LogLine
}

// New instance of log parser service
func NewLogService(dbService DbService, logsConf []LogConf) *logParserService {
	return &logParserService{
		dbService: dbService,
		logsConf:  logsConf,
		LogLines:  make(chan *LogLine),
	}
}

// Run log parser
func (self logParserService) Run() {
	log.Info("LogLine service runs")

	// Gets log lines to chanel self.LogLines
	for _, logConf := range self.logsConf {
		go self.ParseLog(logConf.GetAbsPath(), logConf.GetLogType())
	}


	// wrights log lines from channel to mongo db
	for {
		logLine := <- self.LogLines

		if logLine != nil {
			self.dbService.WrightLog(logLine.Time, logLine.Msg, logLine.Path, logLine.Format)
		}
	}
}


// Parse line by format  " time_in_some_format | msg "
func (self *logParserService) parseLineByFormat(line, path, format string) (*LogLine, error) {
	var (
		timeFormat string
		ok bool
	)

	if timeFormat,ok = dateFormats[format]; !ok {
		msg := fmt.Sprintf("File %s. There is no type format of log %s", path, format)
		return nil, errors.New(msg)
	}

	feilds := strings.Split(line, "|")


	logTime, err := time.Parse(timeFormat, strings.Trim(feilds[0], " "))

	log.Debugf("Parsed time is %s", feilds[0])

	if err != nil || len(feilds) != 2 {
		log.Errorf("Cannot parse line %s in %s type of first_type. " +
			"Format is: %s | This is log message",
			line, path, timeFormat)

		return nil, nil
	}

	msg := strings.Trim(feilds[1], " ")

	return &LogLine{
		Time: logTime,
		Msg: msg,
		Path:path,
		Format: format,
	}, nil
}


// Parse log like tail -f
func (self *logParserService) ParseLog(abspath string, format string) {
	f, err := os.Open(abspath)
	if err != nil {
		log.Errorf("Cannot open file %s. %v", abspath, err)
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	info, err := f.Stat()
	if err != nil {
		log.Errorf("Cannot read file %s. %v", abspath, err)
		return
	}
	oldSize := info.Size()
	for {
		for line, _, err := r.ReadLine(); err != io.EOF;
			line, _, err = r.ReadLine() {
				logLine, err := self.parseLineByFormat(string(line), abspath, format)

				if err != nil {
					log.Error(err.Error())
					return
				}

				self.LogLines <- logLine
		}
		pos, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			panic(err)
		}
		for {
			time.Sleep(time.Second)
			newinfo, err := f.Stat()
			if err != nil {
				panic(err)
			}
			newSize := newinfo.Size()
			if newSize != oldSize {
				if newSize < oldSize {
					f.Seek(0, 0)
				} else {
					f.Seek(pos, io.SeekStart)
				}
				r = bufio.NewReader(f)
				oldSize = newSize
				break
			}
		}
	}
}
