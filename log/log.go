package log

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	logrus "github.com/sirupsen/logrus"
	"io/ioutil"
)

// Logger is the main logger of bioctl
var Logger = logrus.New()

// LoggerBash show [BASH] prefix in logrus message
var LoggerBash = Logger.WithFields(logrus.Fields{
	"prefix": "BASH"})

// SetClassicStyle set logrus.Logger to classic "[2020-01-14 18:53:12] [Level] message" format
func SetClassicStyle(Logger *logrus.Logger) {
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "[2006-01-02 15:04:05]",
	})
}

// SetQuietLog Set quiet log
func SetQuietLog(log *logrus.Logger, quite string) {
	if quite == "true" {
		log.SetOutput(ioutil.Discard)
	}
}

func init() {
	SetClassicStyle(Logger)
}
