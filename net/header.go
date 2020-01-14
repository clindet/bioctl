package net

import "net/http"

// SetDefaultReqHeader set default keep-alive and User-Agent header
func SetDefaultReqHeader(req *http.Request) {
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
}
