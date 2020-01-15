package cmd

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	cio "github.com/openbiox/ligo/io"
	clog "github.com/openbiox/ligo/log"
	"github.com/openbiox/ligo/stringo"
	"github.com/spf13/cobra"
)

var log = clog.Logger
var logBash = clog.LoggerBash
var wd string

func setGlobalFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&(rootClis.Quiet), "quite", "q", "true", "keep slient and drop debug information [true or false]")
	cmd.PersistentFlags().StringVarP(&(rootClis.TaskID), "task-id", "", stringo.RandString(15), "task ID (default is random).")
	cmd.PersistentFlags().StringVarP(&(rootClis.LogDir), "log-dir", "", path.Join(wd, "_log"), "log dir.")
	cmd.PersistentFlags().StringVarP(&(rootClis.SaveLog), "save-log", "", "false", "Save log to file.")
	cmd.PersistentFlags().BoolVarP(&(rootClis.Clean), "clean", "", false, "Remove log dir.")
}
func initCmd() {
	setLog()

	if rootClis.Clean {
		cleanLog()
	}
}

func setLog() {
	rootClis.Quiet = strings.ToLower(rootClis.Quiet)
	rootClis.SaveLog = strings.ToLower(rootClis.SaveLog)
	var logCon io.Writer
	if rootClis.SaveLog == "true" {
		if err := cio.CreateDir(rootClis.LogDir); err != nil {
			return
		}
		var err error
		if logCon, err = cio.Open(fmt.Sprintf("%s/%s.log", rootClis.LogDir, rootClis.TaskID)); err != nil {
			return
		}
	}
	clog.SetLogStream(log, rootClis.Quiet == "true", rootClis.SaveLog == "true", &logCon)
}

func cleanLog() {
	rootClis.HelpFlags = false
	if err := os.RemoveAll(rootClis.LogDir); err != nil {
		log.Warn(err)
	}
}
