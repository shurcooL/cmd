// godoc_router
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	. "gist.github.com/5286084.git"
)

var httpFlag = flag.String("http", ":80", "Listen for HTTP connections on this address")
var gitHubUserFlag = flag.String("github-user", "", "GitHub user with private repos (required)")
var privateGodocHostFlag = flag.String("private-godoc-host", "127.0.0.1:8080", "Host of private Godoc server")

func NewRouter() *httputil.ReverseProxy {
	director := func(r *http.Request) {
		var usePrivate bool
		switch {
		case strings.HasPrefix(r.URL.Path, fmt.Sprintf("/github.com/%s/", *gitHubUserFlag)) ||
			strings.HasPrefix(r.URL.Query().Get("q"), fmt.Sprintf("github.com/%s/", *gitHubUserFlag)):

			usePrivate = true
		case r.URL.Path == "/-/refresh":
			body, err := ioutil.ReadAll(r.Body)
			CheckError(err)
			err = r.Body.Close()
			CheckError(err)
			r.Body = ioutil.NopCloser(bytes.NewReader(body))

			// TODO: Maybe just do it right and use url.ParseQuery?
			usePrivate = strings.HasPrefix(string(body), url.QueryEscape(fmt.Sprintf("path=github.com/%s/", *gitHubUserFlag)))
		case r.URL.Path == "/-/index":
			usePrivate = true
		case r.URL.Path == "/":
			usePrivate = true
		default:
			usePrivate = false
		}

		if usePrivate {
			r.URL.Scheme = "http"
			r.URL.Host = *privateGodocHostFlag
		} else {
			r.URL.Scheme = "http"
			r.URL.Host = "godoc.org"
			r.Host = "godoc.org"
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func main() {
	flag.Parse()

	if *gitHubUserFlag == "" {
		flag.Usage()
		return
	}

	err := http.ListenAndServe(*httpFlag, NewRouter())
	CheckError(err)
}
