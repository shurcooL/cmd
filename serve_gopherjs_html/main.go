package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/shurcooL/go/gopherjs_http"
)

var httpFlag = flag.String("http", ":8080", "Listen for HTTP connections on this address.")

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "usage: serve_gopherjs_html file.html")
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
