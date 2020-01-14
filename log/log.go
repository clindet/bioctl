package log

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	logrus "github.com/sirupsen/logrus"
)

// Logger is the main logger of ganker
var Logger = logrus.New()

func init() {
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "[2006-01-02 15:04:05]",
	})
}
