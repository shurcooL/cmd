// Command gopherjs_serve_html serves an html file with "text/go" script type support.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/shurcooL/go/gopherjs_http"
)

var httpFlag = flag.String("http", ":8080", "Listen for HTTP connections on this address.")

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: gopherjs_serve_html file.html")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		return
	}

	fmt.Printf("serving at http://%s/index.html\n", *httpFlag)

	http.Handle("/index.html", gopherjs_http.HtmlFile(flag.Args()[0]))
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(*httpFlag, nil)
	if err != nil {
		panic(err)
	}
}
