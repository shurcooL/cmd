// Dumps incoming HTTP requests with full detail.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/shurcooL/go-goon"
)

var httpFlag = flag.String("http", ":8080", "Listen for HTTP connections on this address.")

func dumpRequestHandler(w http.ResponseWriter, req *http.Request) {
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(dump))
	if body, err := ioutil.ReadAll(req.Body); err != nil {
		panic(err)
	} else if len(body) <= 64 {
		fmt.Printf("body: %v len: %v\n", body, len(body))
	}
	goon.DumpExpr(req.URL.Query())
	goon.DumpExpr(req.RemoteAddr)
}

func main() {
	flag.Parse()

	fmt.Printf("Starting HTTP request dumper, listening on %q...\n", *httpFlag)

	err := http.ListenAndServe(*httpFlag, http.HandlerFunc(dumpRequestHandler))
	if err != nil {
		panic(err)
	}
}
