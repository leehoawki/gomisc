package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

var m Meta
var f Fields
var local *time.Location
var hostname string
var pid string

func Setup(name string, logstash string) {
	m = Meta{Beat: "logback"}
	f = Fields{Project: name, Service: "golang"}
	local, _ = time.LoadLocation("")
	hostname, _ = os.Hostname()
	pid = strconv.Itoa(os.Getpid())

	logrus.SetFormatter(&logrus.JSONFormatter{})
	conn, err := net.Dial("tcp", logstash)
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(&TcpWriter{
		Cs:         logstash,
		Connection: conn,
	})
	logrus.SetLevel(logrus.InfoLevel)
}

type Fields struct {
	Project string `json:"project"`
	Service string `json:"service"`
}

type Meta struct {
	Beat string `json:"beat"`
}

type TcpWriter struct {
	Cs         string
	Connection net.Conn
}

func (writer *TcpWriter) Write(p []byte) (n int, err error) {
	n, err = writer.Connection.Write(p)
	if err != nil {
		c, e := net.Dial("tcp", writer.Cs)
		if e != nil {
			fmt.Println(string(p))
			return n, err
		} else {
			writer.Connection = c
			n, err = writer.Connection.Write(p)
		}
	}
	return n, err
}

func Debug(message string) {
	logrus.WithFields(fields(message)).Debug()
}

func Info(message string) {
	logrus.WithFields(fields(message)).Info()
}

func Warn(message string) {
	logrus.WithFields(fields(message)).Warn()
}

func Error(err error) {
	logrus.WithFields(fields(err.Error())).Error()
}

func fields(message string) logrus.Fields {
	return logrus.Fields{
		"thread_name": pid,
		"host":        hostname,
		"@timestamp":  time.Now().In(local).Format("2006-01-02T15:04:05.999Z"),
		"logger_name": getCaller(3),
		"@metadata":   m,
		"fields":      f,
		"message":     message,
	}
}

func getCaller(skip int) string {
	_, file, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file
}
