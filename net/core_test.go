package net

import (
	"os"
	"testing"
	"time"

	mpb "github.com/vbauerster/mpb/v4"
)

func TestHttpGetURLs(t *testing.T) {
	urls := []string{"https://raw.githubusercontent.com/openbiox/bioctl/master/net/core.go",
		"https://raw.githubusercontent.com/openbiox/bioctl/master/net/rename.go",
		"https://github.com/Miachol/github_demo", "git@github.com:Miachol/github_demo.git"}
	destDir := []string{os.TempDir(), os.TempDir(), os.TempDir(), os.TempDir() + "/github_demo"}
	param := &Params{}
	param.Retries = 5
	param.Engine = "go-http"
	param.Timeout = 35
	param.Overwrite = false
	param.TaskID = "test"
	param.Thread = 2
	param.LogDir = os.TempDir()
	param.Pbar = mpb.New(
		mpb.WithWidth(45),
		mpb.WithRefreshRate(180*time.Millisecond),
	)
	HTTPGetURLs(urls, destDir, param)
}
