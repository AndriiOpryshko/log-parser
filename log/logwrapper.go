package log

import (
	"fmt"
	"github.com/op/go-logging"
	"net/http"
	"os"
	"io/ioutil"
	"strings"
	"net"
)

var log = logging.MustGetLogger("blitz_api")

// Log format
var format = logging.MustStringFormatter(
	`%{color}time:[%{time:Mon, 02 Jan 2006 15:04:05 MST}] severity:[%{level:.3s}] %{id:03x}%{color:reset} message:[%{message} called_function:[%{callpath}] `,
)

var (
	postfix = ""
	appname = "log-parser"
)

func getExternalIp() string{
	errMsg := "Cannot get external ip. "
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		os.Stderr.WriteString(errMsg + err.Error())
		os.Stderr.WriteString("\n")
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		os.Stderr.WriteString(errMsg + err.Error())
		os.Stderr.WriteString("\n")
		return ""
	}

	var ip string = string(body)

	ip = strings.Trim(ip, "\n")

	return ip
}

func getLocalIp() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func init() {
	postfix = fmt.Sprintf("]\t\t\tservice: [%s], external_ip: [%s], local_ips: [%s]", appname, getExternalIp(), getLocalIp())

	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend1Formatter)

}

func msgPrefix(){
	fmt.Println()
}

// Log for info
func Info(msg string) {
	msgPrefix()
	log.Infof("%s%s", msg,postfix)
}

// Log for notice
func Notice(msg string) {
	msgPrefix()
	log.Noticef("%s%s", msg, postfix)
}

// Log for warning
func Warning(msg string) {
	msgPrefix()
	log.Warning(msg, postfix)
}

// Log for error
func Error(msg string) {
	msgPrefix()
	log.Error(msg, postfix)
}

// Log for critical
func Critical(msg string) {
	msgPrefix()
	log.Critical(msg, postfix)
}

// Log for critical
func Debug(msg string) {
	msgPrefix()
	log.Debug(msg, postfix)
}

// Log for infof
func Infof(format string, args ...interface{}) {
	msgPrefix()
	msg := fmt.Sprintf(format, args...)
	log.Infof("%s%s", msg, postfix)
}

// Log for niticef
func Noticef(format string, args ...interface{}) {
	msgPrefix()
	msg := fmt.Sprintf(format, args...)
	log.Noticef("%s%s", msg, postfix)
}

// Log for warningf
func Warningf(format string, args ...interface{}) {
	msgPrefix()
	msg := fmt.Sprintf(format, args...)
	log.Warningf("%s%s", msg, postfix)
}

// Log for error
func Errorf(format string, args ...interface{}) {
	msgPrefix()
	msg := fmt.Sprintf(format, args...)
	log.Errorf("%s%s", msg, postfix)
}

// Log for critical
func Criticalf(format string, args ...interface{}) {
	msgPrefix()
	msg := fmt.Sprintf(format, args...)
	log.Criticalf("%s%s", msg, postfix)
}

// Log for critical
func Debugf(format string, args ...interface{}) {
	msgPrefix()
	msg := fmt.Sprintf(format, args...)
	log.Debugf("%s%s", msg, postfix)

}