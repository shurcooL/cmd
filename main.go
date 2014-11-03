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
)

var httpFlag = flag.String("http", ":80", "Listen for HTTP connections on this address.")
var gitHubUserFlag = flag.String("github-user", "", "GitHub user with private repos (required).")
var privateGodocHostFlag = flag.String("private-godoc-host", "127.0.0.1:8080", "Host of private Godoc server.")

func NewRouter() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		var usePrivate bool
		switch {
		case strings.HasPrefix(req.URL.Path, fmt.Sprintf("/github.com/%s/", *gitHubUserFlag)) ||
			strings.HasPrefix(req.URL.Query().Get("q"), fmt.Sprintf("github.com/%s/", *gitHubUserFlag)):

			usePrivate = true
		case req.URL.Path == "/-/refresh":
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			err = req.Body.Close()
			if err != nil {
				panic(err)
			}
			req.Body = ioutil.NopCloser(bytes.NewReader(body))

			// TODO: Maybe just do it right and use url.ParseQuery?
			usePrivate = strings.HasPrefix(string(body), "path="+url.QueryEscape(fmt.Sprintf("github.com/%s/", *gitHubUserFlag)))
		case req.URL.Path == "/-/index":
			usePrivate = true
		case req.URL.Path == "/":
			usePrivate = true
		default:
			usePrivate = false
		}

		if usePrivate {
			req.URL.Scheme = "http"
			req.URL.Host = *privateGodocHostFlag
		} else {
			req.URL.Scheme = "http"
			req.URL.Host = "godoc.org"
			req.Host = "godoc.org"
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
	if err != nil {
		panic(err)
	}
}
