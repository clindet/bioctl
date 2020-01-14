package net

import (
	"crypto/tls"
	"math/rand"
	"net"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/openbiox/bioctl/stringo"
)

// SetCollyProxy set colly.Collector with proxy
func SetCollyProxy(c *colly.Collector, pry string, timeout int) {
	trp := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(timeout) * time.Second,
		}).Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	proxyPool := []string{}
	if pry == "" {
		return
	}
	if strings.Contains(pry, ";") {
		proxyPool = stringo.StrSplit(pry, ";", 1000000)
	} else if pry != "" {
		proxyPool = append(proxyPool, pry)
	}
	urlproxy, _ := RandProxy(pry)
	trp.Proxy = http.ProxyURL(urlproxy)

	if len(proxyPool) > 1 {
		urlproxy, _ := RandProxy(pry)
		trp.Proxy = http.ProxyURL(urlproxy)
	}
	c.WithTransport(trp)
}

// RandProxy return a proxy from proxy string
func RandProxy(proxy string) (*neturl.URL, string) {
	if proxy == "" {
		return nil, ""
	}
	proxyPool := []string{}
	if strings.Contains(proxy, ";") {
		proxyPool = stringo.StrSplit(proxy, ";", 1000000)
	} else {
		proxyPool = append(proxyPool, proxy)
	}
	i := rand.Int63n(int64(len(proxyPool) - 0))
	urli := neturl.URL{}
	urlproxy, _ := urli.Parse(proxyPool[i])
	return urlproxy, proxyPool[i]
}
