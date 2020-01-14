package par

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	bexec "github.com/openbiox/bioctl/exec"
	cio "github.com/openbiox/bioctl/io"
	clog "github.com/openbiox/bioctl/log"
	"github.com/openbiox/bioctl/stringo"
	"github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
)

// ClisT is the type of parameters of parTasks
type ClisT struct {
	Script      string
	Index       string
	Env         string
	ForceAddIdx string
	Thread      int
	TaskID      string
	LogDir      string
	Quiet       string
}

var log = clog.Logger

// Tasks parallel run tasks
func Tasks(parClis *ClisT) error {
	index2 := []int{}
	clog.SetQuietLog(log, parClis.Quiet)
	if parClis.Index != "" {
		index := strings.Split(parClis.Index, ",")
		for i := range index {
			if strings.Contains(index[i], "-") {
				startEnd := strings.Split(index[i], "-")
				start, _ := strconv.ParseInt(startEnd[0], 10, 64)
				end, _ := strconv.ParseInt(startEnd[1], 10, 64)
				for j := start; j < end+1; j++ {
					index2 = append(index2, int(j))
				}
			} else {
				val, _ := strconv.ParseInt(index[i], 10, 64)
				index2 = append(index2, int(val))
			}
		}
	} else {
		index2 = append(index2, 1)
	}
	envObj := make(map[string]string)
	if parClis.Env != "" {
		envSlice := strings.Split(parClis.Env, ",")
		for k := range envSlice {
			envSlice2 := strings.Split(envSlice[k], ":")
			envObj[envSlice2[0]] = envSlice2[1]
		}
	}
	sort.Sort(sort.IntSlice(index2))

	sem := make(chan bool, parClis.Thread)
	var wg sync.WaitGroup
	p := mpb.New(
		mpb.WithWaitGroup(&wg),
		mpb.WithWidth(60),
	)
	wg.Add(len(index2))
	err := cio.CreateDir(parClis.LogDir)
	if err != nil {
		return err
	}
	logSlice := []string{}
	for i := range index2 {
		logSlice = append(logSlice, fmt.Sprintf("%s/%s-%d.log", parClis.LogDir, parClis.TaskID, i+1))
	}
	log.Infof("Total %d tasks were submited (%s | %s).", len(index2), parClis.TaskID, parClis.Index)
	log.Infof("Timestamp: %s", time.Now())
	hostname, _ := os.Hostname()
	user, _ := user.Current()
	platform := runtime.GOOS
	log.Infof("Hostname: %s, Username: %s", hostname, user.Username)
	wd, _ := os.Getwd()
	log.Infof("Platform: %s, Working: %s", platform, wd)
	log.Infof("Task log from #1 to #%d will be saved in %s/%s-*.log .", len(index2), parClis.LogDir, parClis.TaskID)
	errorMsg := make(map[int]string)
	var lock sync.Mutex
	total := 100
	for i := 0; i < len(index2); i++ {
		var ind = index2[i]
		sem <- true
		name := fmt.Sprintf("Task #%d", i+1)
		bar := p.AddBar(int64(total), mpb.BarID(i),
			mpb.BarStyle("╢=>-╟"),
			// override mpb.DefaultSpinnerStyle
			mpb.PrependDecorators(
				// simple name decorator
				decor.OnComplete(decor.Spinner(nil, decor.WCSyncSpace), "√"),
				decor.Name(" "+name+fmt.Sprintf(" | index: %5d#", ind)),
			),
			mpb.AppendDecorators(
				// replace ETA decorator with "done" message, OnComplete event
				decor.Name(" | elapsed: "), decor.Elapsed(decor.ET_STYLE_HHMMSS),
			),
		)
		// simulating some work
		go func(i int, bar *mpb.Bar) {
			defer func() {
				<-sem
			}()
			var logPath = logSlice[i]
			defer wg.Done()
			go func(bar *mpb.Bar) {
				count := 1
				for {
					if bar.Completed() {
						break
					}
					bar.SetCurrent(int64(count))
					count++
					if count == 99 {
						count = 0
					}
					time.Sleep(time.Second / 5)
				}
			}(bar)
			var cmd *exec.Cmd
			script := stringo.StrReplaceAll(parClis.Script, "\n$", "")

			for k, v := range envObj {
				script = stringo.StrReplaceAll(script, fmt.Sprintf("{{%s}}", k), v)
			}
			indStr := fmt.Sprintf("%d", ind)
			if strings.Contains(script, "{{index}}") {
				script = stringo.StrReplaceAll(script, "{{index}}", indStr)
			} else if parClis.ForceAddIdx == "true" {
				script = script + " " + indStr
			}

			cmd = exec.Command("bash", "-c", script)
			err := bexec.System(cmd, logPath, true)
			if err != nil {
				//fmt.Println(err)
				lock.Lock()
				errorMsg[i+1] = fmt.Sprintf("Task #%d error: %s", i+1, err)
				lock.Unlock()
				bar.SetCurrent(int64(0))
				bar.Abort(false)
				time.Sleep(time.Second * 1)
			} else {
				bar.SetCurrent(int64(total))
			}
		}(i, bar)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	// wait for all bars to complete and flush
	p.Wait()
	for i := 0; i < len(index2); i++ {
		if errorMsg[i+1] != "" {
			log.Warn(errorMsg[i+1])
		}
	}
	return nil
}
