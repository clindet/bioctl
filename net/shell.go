package net

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	cio "github.com/openbiox/bioctl/io"
)

func setWgetCmds(url string, destFn string, opt *Params) *exec.Cmd {
	cmdCheck := exec.Command("sh", "-c", `wget --help | grep '\--show-progress'`)
	isNewWget, _ := cmdCheck.CombinedOutput()
	args := []string{"-c", url, "-O", destFn,
		"--connect-timeout=" + strconv.Itoa(opt.Timeout)}
	if string(isNewWget) != "" {
		args = append(args, "--show-progress")
	} else {
		args = append(args, "--progress=bar")
	}
	if opt.ExtraArgs != "" {
		extraArgsList := strings.Split(opt.ExtraArgs, " ")
		args = append(args, extraArgsList...)
	}
	cmd := exec.Command("wget", args...)
	SetCmdProxyEnv(cmd, opt.Proxy, url)
	return cmd
}

// Wget use wget to download files
func Wget(url string, destFn string, opt *Params) (err error) {
	if url == "" {
		return errors.New("at least one of URL is required")
	}
	cmd := setWgetCmds(url, destFn, opt)
	logPath := ""
	if opt.SaveLog {
		logPath = path.Join(opt.LogDir, fmt.Sprintf("%s_%s_wget.log",
			opt.TaskID, path.Base(destFn)))
		log.Infof("Download Log of %s saved to => %s", url, logPath)
		cio.CreateFileParDir(logPath)
	}
	time.Sleep(1 * time.Second)
	err = RetriesURL(url, cmd, logPath, opt)
	if err != nil {
		return err
	}
	return nil
}

// Curl use curl to download files
func Curl(url string, destFn string, opt *Params) (err error) {
	if url == "" {
		return errors.New("at least one of URL is required")
	}
	args := []string{url, "-o", destFn, "--connect-timeout", strconv.Itoa(opt.Timeout)}
	if opt.ExtraArgs != "" {
		extraArgsList := strings.Split(opt.ExtraArgs, " ")
		args = append(args, extraArgsList...)
	}
	cmd := exec.Command("curl", args...)
	SetCmdProxyEnv(cmd, opt.Proxy, url)
	logPath := ""
	if opt.SaveLog {
		logPath = path.Join(opt.LogDir, fmt.Sprintf("%s_%s_curl.log", opt.TaskID,
			path.Base(destFn)))
		log.Infof("Download Log of %s saved to => %s", url, logPath)
		cio.CreateFileParDir(logPath)
	}
	time.Sleep(1 * time.Second)
	err = RetriesURL(url, cmd, logPath, opt)
	if err != nil {
		return err
	}
	return nil
}

func setAxelCmd(url string, destFn string, opt *Params) *exec.Cmd {
	cmdCheck := exec.Command("sh", "-c", `axel --help | grep '\--timeout'`)
	isNewAxel, _ := cmdCheck.CombinedOutput()
	args := []string{url, "-o", destFn, "-n", strconv.Itoa(opt.AxelThread)}
	if string(isNewAxel) != "" {
		args = append(args, []string{"--timeout=" + strconv.Itoa(opt.Timeout)}...)
	}
	if opt.ExtraArgs != "" {
		extraArgsList := strings.Split(opt.ExtraArgs, " ")
		args = append(args, extraArgsList...)
	}
	cmd := exec.Command("axel", args...)
	return cmd
}

// Axel use axel to download files
func Axel(url string, destFn string, opt *Params) (err error) {
	if url == "" {
		return errors.New("at least one of URL is required")
	}
	cmd := setAxelCmd(url, destFn, opt)
	SetCmdProxyEnv(cmd, opt.Proxy, url)
	logPath := ""
	if opt.SaveLog {
		logPath = path.Join(opt.LogDir,
			fmt.Sprintf("%s_%s_axel.log", opt.TaskID, path.Base(destFn)))
		log.Infof("Download Log of %s saved to => %s", url, logPath)
		cio.CreateFileParDir(logPath)
	}
	time.Sleep(1 * time.Second)
	err = RetriesURL(url, cmd, logPath, opt)
	if err != nil {
		return err
	}
	return nil
}

// Git use git to download files
func Git(url string, destFn string, opt *Params) (err error) {
	if url == "" {
		return errors.New("at least one of URL is required")
	}
	args := []string{"clone", "--recursive", "--progress"}
	if opt.ExtraArgs != "" {
		extraArgsList := strings.Split(opt.ExtraArgs, " ")
		args = append(args, extraArgsList...)
	}
	args = append(args, url, destFn)
	cmd := exec.Command("git", args...)
	SetCmdProxyEnv(cmd, opt.Proxy, url)
	logPath := ""
	if opt.SaveLog {
		logPath = path.Join(opt.LogDir, fmt.Sprintf("%s_%s_git.log",
			opt.TaskID, path.Base(destFn)))
		log.Infof("Download Log of %s saved to => %s", url, logPath)
		cio.CreateFileParDir(logPath)
	}
	time.Sleep(1 * time.Second)
	err = RetriesURL(url, cmd, logPath, opt)
	if err != nil {
		return err
	}
	return nil
}

// Hg use hg to download files
func Hg(url string, destFn string, opt *Params) (err error) {
	if url == "" {
		return errors.New("at least one of URL is required")
	}
	args := []string{"clone"}
	if opt.ExtraArgs != "" {
		extraArgsList := strings.Split(opt.ExtraArgs, " ")
		args = append(args, extraArgsList...)
	}
	args = append(args, url, destFn)
	cmd := exec.Command("hg", args...)
	SetCmdProxyEnv(cmd, opt.Proxy, url)
	logPath := ""
	if opt.SaveLog {
		logPath = path.Join(opt.LogDir, fmt.Sprintf("%s_%s_hg.log",
			opt.TaskID, path.Base(destFn)))
		log.Infof("Download Log of %s saved to => %s", url, logPath)
		cio.CreateFileParDir(logPath)
	}
	time.Sleep(1 * time.Second)
	oldRetries := opt.Retries
	opt.Retries = 1
	err = RetriesURL(url, cmd, logPath, opt)
	opt.Retries = oldRetries
	if err != nil {
		cmd := exec.Command("git", args...)
		err = RetriesURL(url, cmd, logPath, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

// Rsync use rsync to download files
func Rsync(url string, destFn string, opt *Params) (err error) {
	if url == "" {
		return errors.New("at least one of URL is required")
	}
	args := []string{url, destFn}
	if opt.ExtraArgs != "" {
		extraArgsList := strings.Split(opt.ExtraArgs, " ")
		args = append(args, extraArgsList...)
	}
	cmd := exec.Command("rsync", args...)
	SetCmdProxyEnv(cmd, opt.Proxy, url)
	logPath := ""
	if opt.SaveLog {
		logPath = path.Join(opt.LogDir, fmt.Sprintf("%s_%s_rsync.log", opt.TaskID, path.Base(destFn)))
		log.Infof("Download Log of %s saved to => %s", url, logPath)
		cio.CreateFileParDir(logPath)
	}
	time.Sleep(1 * time.Second)
	err = RetriesURL(url, cmd, logPath, opt)
	if err != nil {
		return err
	}
	return nil
}

// GdcClient use gdc-client to download files
func GdcClient(fileID string, manifest string, outDir string, opt *Params) (err error) {
	if fileID == "" && manifest == "" {
		return errors.New("at least one of fileID or manifest is required")
	}
	args := []string{}
	if manifest == "" {
		args = []string{"download", fileID, "-d", outDir}
	} else {
		args = []string{"download", "-m", manifest, "-d", outDir}
	}
	if opt.ExtraArgs != "" {
		extraArgsList := strings.Split(opt.ExtraArgs, " ")
		args = append(args, extraArgsList...)
	}
	if opt.Token != "" {
		args = append(args, "-t", opt.Token)
	}
	cmd := exec.Command("gdc-client", args...)
	SetCmdProxyEnv(cmd, opt.Proxy, "")
	taskName := ""
	if fileID != "" {
		taskName = fileID
	} else {
		taskName = manifest
	}
	logPath := ""
	if opt.SaveLog {
		logPath = path.Join(opt.LogDir, fmt.Sprintf("%s_gdc-client.log", opt.TaskID))
		log.Infof("Download Log of %s saved to => %s", taskName, logPath)
		cio.CreateFileParDir(logPath)
	}
	time.Sleep(1 * time.Second)
	err = RetriesTask(taskName, cmd, logPath, opt)
	if err != nil {
		return err
	}
	return nil
}

// Prefetch use sra-tools prefetch to download files
func Prefetch(srr string, krt string, outDir string, opt *Params) (err error) {
	if srr == "" && krt == "" {
		return errors.New("at least one of srr or krt is required")
	}
	args := []string{"-O", outDir, "-X", "500GB"}
	if opt.ExtraArgs != "" {
		extraArgsList := strings.Split(opt.ExtraArgs, " ")
		args = append(args, extraArgsList...)
	}
	if krt == "" {
		args = append(args, srr)
	} else {
		args = append(args, krt)
	}

	cmd := exec.Command("prefetch", args...)
	SetCmdProxyEnv(cmd, opt.Proxy, "")
	logPath := ""
	if opt.SaveLog {
		logPath = path.Join(opt.LogDir, fmt.Sprintf("%s_prefetch.log", opt.TaskID))
	}
	taskName := ""
	if srr != "" {
		taskName = srr
	} else {
		taskName = krt
	}
	if opt.SaveLog {
		log.Infof("Download Log of %s saved to => %s", taskName, logPath)
		cio.CreateFileParDir(logPath)
	}
	time.Sleep(1 * time.Second)
	err = RetriesTask(taskName, cmd, logPath, opt)
	if err != nil {
		return err
	}
	return nil
}

// Egafetch use pyega3 fetch EGA Archive files
// https://ega-archive.org/download/downloader-quickguide-v3#defineCredentials
func Egafetch(ega string, fileID string, outDir string, opt *Params) (err error) {
	if ega == "" && fileID == "" {
		return errors.New("at least one of ega or fileID is required")
	}
	args := []string{"fetch", "--saveto", outDir, "-cf", opt.EgaCredentials}
	if opt.ExtraArgs != "" {
		extraArgsList := strings.Split(opt.ExtraArgs, " ")
		args = append(args, extraArgsList...)
	} else {
		args = append(args, []string{"--format", "BAM"}...)
	}
	if fileID == "" {
		args = append(args, ega)
	} else {
		args = append(args, fileID)
	}

	cmd := exec.Command("pyega3", args...)
	SetCmdProxyEnv(cmd, opt.Proxy, "")
	logPath := ""
	if opt.SaveLog {
		logPath = path.Join(opt.LogDir, fmt.Sprintf("%s_pyega3.log", opt.TaskID))
	}
	taskName := ""
	if ega != "" {
		taskName = ega
	} else {
		taskName = fileID
	}
	if opt.SaveLog {
		log.Infof("Download Log of %s saved to => %s", taskName, logPath)
		cio.CreateFileParDir(logPath)
	}
	time.Sleep(1 * time.Second)
	err = RetriesTask(taskName, cmd, logPath, opt)
	if err != nil {
		return err
	}
	return nil
}
