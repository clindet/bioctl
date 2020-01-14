package net

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	cio "github.com/openbiox/bioctl/io"
	"github.com/openbiox/bioctl/stringo"
	mpb "github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
)

func checkHTTPGetURLRdirect(resp *http.Response, url string, destFn string, opt *Params) (status bool) {
	if strings.Contains(url, "https://www.sciencedirect.com") {
		v, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			if stringo.StrDetect(string(v), "https://pdf.sciencedirectassets.com") {
				url = stringo.StrExtract(string(v), `https://pdf.sciencedirectassets.com/.*&type=client`, 1)[0]
				HTTPGetURL(url, destFn, opt)
				return true
			}
		}
	}
	return false
}

func defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 20 {
		return errors.New("stopped after 20 redirects")
	}
	return nil
}

func checkResp(resp *http.Response) (err error) {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("access failed: %s", resp.Request.URL.String())
	}
	return nil
}

func checkGitEngine(url string) string {
	if stringo.StrDetect(url, "^git@") {
		return "git"
	}
	hasHg := cio.CheckInputFiles([]string{"hg"}) == nil
	sitesGit := []string{"github.com", "gitlab.com", "bitbucket.org"}
	sitesHg := []string{"bitbucket.org"}
	url = stringo.StrReplaceAll(url, "/$", "")
	for _, v := range sitesHg {
		if stringo.StrDetect(url, v) && strings.Count(url, "/") <= 4 && hasHg {
			return "hg"
		}
	}
	for _, v := range sitesGit {
		if stringo.StrDetect(url, v) && strings.Count(url, "/") <= 4 {
			return "git"
		}
	}
	return ""
}

func downloadWorker(client *http.Client, req *http.Request, url string, destFn string, opt *Params) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if !opt.Quiet {
			log.Warnf("Access failed: %s", url)
			fmt.Println("")
		}
		return err
	}
	if checkHTTPGetURLRdirect(resp, url, destFn, opt) {
		//return errors.New("Redirect fail")
		return nil
	}
	size := resp.ContentLength

	if hasParDir, _ := cio.PathExists(filepath.Dir(destFn)); !hasParDir {
		err = cio.CreateFileParDir(destFn)
		if err != nil {
			return err
		}
	}
	// create dest
	destName := filepath.Base(url)
	dest, err := os.Create(destFn)
	if err != nil {
		log.Warnf("Can't create %s: %v\n", destName, err)
		return err
	}
	defer dest.Close()
	prefixStr := filepath.Base(destFn)
	prefixStrLen := utf8.RuneCountInString(prefixStr)
	if prefixStrLen > 35 {
		prefixStr = prefixStr[0:31] + "..."
	}
	prefixStr = fmt.Sprintf("%-35s\t", prefixStr)
	if !opt.Quiet {
		bar := opt.Pbar.AddBar(size,
			mpb.BarStyle("[=>-|"),
			mpb.PrependDecorators(
				decor.Name(prefixStr, decor.WC{W: len(prefixStr) + 1, C: decor.DidentRight}),
				decor.CountersKibiByte("% -.1f / % -.1f\t"),
				decor.OnComplete(decor.Percentage(decor.WC{W: 5}), " "+"√"),
			),
			mpb.AppendDecorators(
				decor.EwmaETA(decor.ET_STYLE_MMSS, float64(size)/2048),
				decor.Name(" ] "),
				decor.AverageSpeed(decor.UnitKiB, "% .1f"),
			),
		)
		// create proxy reader
		reader := bar.ProxyReader(resp.Body)
		// and copy from reader, ignoring errors
		_, err = io.Copy(dest, reader)
		if err != nil {
			reader.Close()
			bar.Abort(true)
			return err
		}
	} else {
		_, err = io.Copy(dest, io.Reader(resp.Body))
		if err != nil {
			return err
		}
	}
	return nil
}
