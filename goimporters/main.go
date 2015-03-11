// Displays an import graph of Go packages that import the specified Go package in your GOPATH workspace.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/kisielk/gotool"
	"github.com/shurcooL/go-goon"
	"github.com/shurcooL/go/u/u3"
	"golang.org/x/tools/refactor/importgraph"
)

func init() {
	if _, err := exec.LookPath("dot"); err != nil {
		// TODO: Replace dot with an importable native Go package to get rid of this annoying external dependency.
		fmt.Fprintln(os.Stderr, "`dot` command is required (try `brew install graphviz` to install it).")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: goimporters package")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	importPathPatterns := flag.Args()
	importPaths := gotool.ImportPaths(importPathPatterns)
	if len(importPaths) != 1 {
		fmt.Fprintf(os.Stderr, "need to specify 1 package, but matched %v\n", len(importPaths))
		flag.Usage()
		os.Exit(2)
		return
	}
	target := importPaths[0]

	forward, reverse, graphErrors := importgraph.Build(&build.Default)
	_, _, _ = forward, reverse, graphErrors
	if graphErrors != nil {
		goon.DumpExpr(graphErrors)
		panic(0)
	}

	reachable := reverse.Search(target)

	renderGraph := func() ([]byte, error) {
		var in bytes.Buffer

		fmt.Fprintln(&in, `digraph "" {`)
		for k, v := range forward {
			for k2 := range v {
				if !reachable[k] || !reachable[k2] {
					continue
				}

				fmt.Fprintf(&in, "	%q -> %q;\n", k, k2)
			}
		}
		in.WriteString("}")

		cmd := exec.Command("dot", "-Tsvg")
		cmd.Stdin = &in
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		if i := bytes.Index(out, []byte("<svg")); i < 0 {
			return nil, errors.New("<svg not found")
		} else {
			out = out[i:]
		}
		return out, nil
	}

	graphSvg, err := renderGraph()
	if err != nil {
		log.Panicln("renderGraph:", err)
	}

	mux := http.NewServeMux()
	stopServerChan := make(chan struct{})
	mux.HandleFunc("/index", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "<html><body>")
		w.Write(graphSvg)
		io.WriteString(w, "</body></html>")

		go func() {
			time.Sleep(time.Second)
			stopServerChan <- struct{}{}
		}()
	})
	u3.DisplayHtmlInBrowser(mux, stopServerChan, "")
}
