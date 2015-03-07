// Command gopherjs_serve_html serves an html file with "text/go" script type support.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/shurcooL/go/gists/gist8065433"
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

	// Print all addresses that are being served.
	var hosts []string
	if len(*httpFlag) >= 1 && (*httpFlag)[0] == ':' {
		ips, err := gist8065433.GetAllIps()
		if err != nil {
			panic(err)
		}
		for _, ip := range ips {
			if ip == "127.0.0.1" {
				ip = "localhost"
			}
			hosts = append(hosts, ip+*httpFlag)
		}
	} else {
		hosts = []string{*httpFlag}
	}
	fmt.Println("serving, available at:")
	for _, host := range hosts {
		fmt.Printf("http://%s/index.html\n", host)
	}

	http.Handle("/index.html", gopherjs_http.HtmlFile(flag.Args()[0]))
	http.Handle("/", http.FileServer(http.Dir("./")))

	err := http.ListenAndServe(*httpFlag, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
