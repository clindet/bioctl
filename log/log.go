package log

import (
	"io/ioutil"
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	logrus "github.com/sirupsen/logrus"
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

// New Creates a new logger.
func New() *logrus.Logger {
	return logrus.New()
}

// SetQuietLog Set quiet log
func SetQuietLog(log *logrus.Logger, quite string) {
	if quite == "true" {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(os.Stderr)
	}
}

func init() {
	SetClassicStyle(Logger)
}
