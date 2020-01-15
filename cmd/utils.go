package cmd

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	cvrt "github.com/openbiox/ligo/convert"
	cio "github.com/openbiox/ligo/io"
	clog "github.com/openbiox/ligo/log"
	"github.com/openbiox/ligo/stringo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = clog.Logger
var logBash = clog.LoggerBash
var logEnv = log.WithFields(logrus.Fields{
	"prefix": "Env"})
var wd string

func setGlobalFlag(cmd *cobra.Command) {
	wd, _ = os.Getwd()
	cmd.PersistentFlags().BoolVarP(&(rootClis.Quiet), "quite", "q", false, "no output")
	cmd.PersistentFlags().StringVarP(&(rootClis.TaskID), "task-id", "", stringo.RandString(15), "task ID (default is random).")
	cmd.PersistentFlags().StringVarP(&(rootClis.LogDir), "log-dir", "", path.Join(wd, "_log"), "log dir.")
	cmd.PersistentFlags().BoolVarP(&(rootClis.SaveLog), "save-log", "s", false, "Save log to file.")
	cmd.PersistentFlags().BoolVarP(&(rootClis.Clean), "clean", "", false, "Remove log dir.")
}
func initCmd(cmd *cobra.Command, args []string) {
	setLog()
	logEnv.Infof("prog: %s", cmd.CommandPath())
	logEnv.Infof("args: %s", strings.Join(args, " "))
	logEnv.Infof("env (global): %v", cvrt.Struct2Map(rootClis))
	if rootClis.Clean {
		cleanLog()
	}
}

func setLog() {
	var logCon io.Writer
	var logDir = rootClis.LogDir
	var logPrefix string

	if rootClis.SaveLog {
		if logDir == "" {
			logDir = filepath.Join(os.TempDir(), "_log")
		}
		logPrefix = fmt.Sprintf("%s/%s", logDir, rootClis.TaskID)
		cio.CreateDir(logDir)
		logCon, _ = cio.Open(logPrefix + ".log")
	}
	clog.SetLogStream(log, rootClis.Quiet, rootClis.SaveLog, &logCon)
}

func cleanLog() {
	rootClis.HelpFlags = false
	if err := os.RemoveAll(rootClis.LogDir); err != nil {
		log.Warn(err)
	}
}
