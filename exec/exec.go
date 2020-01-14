package exec

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	cio "github.com/openbiox/bioctl/io"
	clog "github.com/openbiox/bioctl/log"
	"github.com/sirupsen/logrus"
)

var stdout, stderr []byte
var errStdout, errStderr error

// SystemQuiet run cmd using exec.Cmd object (slient mode)
func SystemQuiet(cmd *exec.Cmd, logPath string) error {
	return System(cmd, logPath, true)
}

// System using exec.Cmd object (can control wheature output to console)
func System(cmd *exec.Cmd, logPath string, quiet bool) error {
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var log = clog.New()
	var logBash = log.WithFields(logrus.Fields{
		"prefix": "BASH"})

	cmdStr := ""
	if logPath != "/dev/null" && logPath != "" {
		cmdStr = strings.Join(cmd.Args, " ") + " &>> " + logPath
	} else {
		cmdStr = strings.Join(cmd.Args, " ")
	}
	var logCon *os.File
	var err error
	if logPath != "" {
		logCon, err = cio.Open(logPath)
		if err != nil {
			return err
		}
		if quiet {
			log.SetOutput(io.Writer(logCon))
		} else {
			log.SetOutput(io.MultiWriter(os.Stderr, logCon))
		}
	} else if quiet {
		log.SetOutput(ioutil.Discard)
	}
	logBash.Info(cmdStr)
	err1 := cmd.Start()
	go func() {
		if stdoutIn != nil {
			stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn, quiet, logCon)
		}
	}()
	go func() {
		if stderrIn != nil {
			stderr, errStderr = copyAndCapture(os.Stderr, stderrIn, quiet, logCon)
		}
	}()
	err2 := cmd.Wait()
	defer logCon.Close()
	errs := []error{err1, err2}
	for i := range errs {
		if errs[i] != nil {
			log.Warningf("cmd.Run() failed with %s\n", errs[i])
			log.SetOutput(os.Stderr)
			return errs[i]
		}
	}

	log.SetOutput(os.Stderr)
	return nil
}

func copyAndCapture(w io.Writer, r io.Reader, quiet bool, con *os.File) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			if !quiet {
				w.Write(d)
			}
			if con != nil {
				con.Write(d)
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}
