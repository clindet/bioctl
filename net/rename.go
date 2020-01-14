package net

import (
	"mime"
	"net/http"
	neturl "net/url"
	"path"
	"strings"

	"github.com/openbiox/bioctl/stringo"
)

func FormatURLfileName(url string, remoteName bool, timeout int, proxy string) (fname string) {
	if stringo.StrDetect(url, "^git@") {
		return path.Base(url)
	}
	if remoteName && !strings.Contains(url, "ndownloader.figshare.com") {
		client := NewHTTPClient(timeout, proxy)
		req, err := http.NewRequest("GET", url, nil)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		return resp2Filname(resp)
	} else if remoteName {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		resp, err := client.Head(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		return resp2Filname(resp)
	}
	u, _ := neturl.Parse(url)
	uQ := u.Query()
	fname = path.Base(url)
	// cell.com
	if stringo.StrDetect(url, "/pdfExtended/|/pdfdirect/|/Article/Pdf/|/content/articlepdf/|/rmp/pdf/") {
		fname = path.Base(url) + ".pdf"
	} else if stringo.StrDetect(url, "showPdf[?]pii=") {
		fname = path.Base(stringo.StrReplaceAll(url, "showPdf[?]pii=", "")) + ".pdf"
	} else if stringo.StrDetect(url, "track/pdf") {
		fname = path.Base(url) + ".pdf"
	} else if stringo.StrDetect(url, "&type=printable") {
		fname = strings.ReplaceAll(path.Base(url), "&type=printable", "") + ".pdf"
	} else if fname == "pdf" {
		fname = path.Base(strings.ReplaceAll(url, "/pdf", ".pdf"))
	} else if stringo.StrDetect(fname, "[?]Expires=") {
		fname = stringo.StrReplaceAll(fname, "[?]Expires=.*", "")
	} else if stringo.StrDetect(url, "/action/downloadSupplement[?].*") {
		fname = stringo.StrReplaceAll(fname, "downloadSupplement.*file=", "")
	} else if stringo.StrDetect(url, "(.com/doi/pdf/)|(.org/doi/pdf/)|(.org/doi/pdfdirect/)") {
		if stringo.StrDetect(url, "[?]articleTools=true") {
			fname = stringo.StrReplaceAll(fname, "[?]articleTools=true", "")
		}
		fname = fname + ".pdf"
	} else if stringo.StrDetect(url, "[?]md5=.*&pid=.*") {
		fname = stringo.StrReplaceAll(fname, "[?]md5=.*&pid=", "")
	} else if stringo.StrDetect(fname, "[?]download=true$") {
		fname = stringo.StrReplaceAll(fname, "[?]download=true$", "")
	} else if stringo.StrDetect(fname, "[?]_hash=.*") {
		fname = stringo.StrReplaceAll(fname, "[?]_hash=.*", "")
	} else if stringo.StrDetect(url, "sd/pdf/render") {
		fname = "supp." + fname + ".pdf"
	} else if strings.Contains(url, "https://www.ncbi.nlm.nih.gov/geo/download/?acc=") {
		if strings.Contains(url, "file&file=") {
			fname = uQ["file"][0]
		} else {
			fname = uQ["acc"][0] + ".tar"
		}
		fname, _ = neturl.QueryUnescape(fname)
	} else if strings.Contains(url, "www.ncbi.nlm.nih.gov/geo/query/acc") {
		fname = uQ["acc"][0] + ".txt"
	} else if strings.Contains(url, "https://www.ncbi.nlm.nih.gov/Traces/study/backends") &&
		strings.Contains(url, "rt_all") &&
		strings.Contains(url, "rs=") {
		fname = stringo.StrReplaceAll(uQ["rs"][0], `[(]primary_search_ids:|[)]|"`, "") + "_SraRunTable.txt"
	} else if strings.Contains(url, "https://www.ncbi.nlm.nih.gov/Traces/study/backends") &&
		strings.Contains(url, "acc_all") &&
		strings.Contains(url, "rs=") {
		fname = stringo.StrReplaceAll(uQ["rs"][0], `[(]primary_search_ids:|[)]|"`, "") + "_SraAccList.txt"
	}
	if strings.Contains(fname, "needAccess=true") {
		fname = stringo.StrReplaceAll(fname, "[?]needAccess=true", "")
	}
	return fname
}

func resp2Filname(resp *http.Response) (fname string) {
	contentDis := resp.Header.Get("Content-Disposition")
	if contentDis != "" && strings.Contains(contentDis, "filename") {
		_, params, err := mime.ParseMediaType(contentDis)
		if err != nil {
			log.Warn(err)
		} else {
			fname = params["filename"]
		}
	}
	return fname
}
