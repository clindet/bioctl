package net

import (
	"net/http"
	"net/http/cookiejar"
	"sync"

	clog "github.com/openbiox/bioctl/log"
	mpb "github.com/vbauerster/mpb/v4"
)

var gCurCookies []*http.Cookie
var gCurCookieJar *cookiejar.Jar

// Params is the type object to run bioctl net function
type Params struct {
	// go-http, wget, curl, axel,
	Engine         string
	Token          string
	EgaCredentials string
	Mirror         string
	Thread         int
	AxelThread     int
	ExtraArgs      string
	Proxy          string
	TaskID         string
	Quiet          bool
	Overwrite      bool
	Ignore         bool
	LogDir         string
	SaveLog        bool
	Retries        int
	Timeout        int
	RetSleepTime   int
	RemoteName     bool
	Pbar           *mpb.Progress
}

var lock sync.Mutex
var log = clog.Logger
