package net

import (
	"os/exec"
	"strings"
)

// SetCmdProxyEnv set downloader (shell tools) env
func SetCmdProxyEnv(cmd *exec.Cmd, proxy, url string) {
	if proxy == "" {
		return
	}
	proxyFlag := ""
	if strings.Contains(url, "https://") {
		proxyFlag = "https_proxy=" + proxy
	} else if strings.Contains(url, "ftp://") {
		proxyFlag = "ftp_proxy=" + proxy
	} else {
		proxyFlag = "http_proxy=" + proxy
	}
	_, proxy = RandProxy(proxy)
	cmd.Env = append(cmd.Env, proxyFlag)
	cmd.Env = append(cmd.Env, "use_proxy=on")
}
