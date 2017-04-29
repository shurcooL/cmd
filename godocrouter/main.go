// godocrouter is a reverse proxy that augments a private godoc server instance with global godoc.org instance.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var (
	httpFlag             = flag.String("http", ":80", "Listen for HTTP connections on this address.")
	githubUserFlag       = flag.String("github-user", "", "GitHub user with private repos (required).")
	privateGodocHostFlag = flag.String("private-godoc-host", "127.0.0.1:8080", "Host of private Godoc server.")
)

func main() {
	flag.Parse()

	if *githubUserFlag == "" {
		flag.Usage()
		os.Exit(2)
	}

	err := http.ListenAndServe(*httpFlag, newRouter())
	if err != nil {
		log.Fatalln(err)
	}
}

func newRouter() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		var usePrivate bool
		switch {
		case strings.HasPrefix(req.URL.Path, fmt.Sprintf("/github.com/%s/", *githubUserFlag)) ||
			strings.HasPrefix(req.URL.Query().Get("q"), fmt.Sprintf("github.com/%s/", *githubUserFlag)):

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
			usePrivate = strings.HasPrefix(string(body), "path="+url.QueryEscape(fmt.Sprintf("github.com/%s/", *githubUserFlag)))
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
