package cmd

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func setQuietLog(log *logrus.Logger, quite string) {
	if rootClis.Quite == "true" {
		log.SetOutput(ioutil.Discard)
	}
}
