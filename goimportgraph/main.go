// Displays an import graph within specified Go packages.
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
	"github.com/shurcooL/go/u/u3"
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
	fmt.Fprint(os.Stderr, "Usage: goimportgraph packages\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, `
Examples:
  goimportgraph encoding/...

  goimportgraph github.com/shurcooL/Conception-go/...
`)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	importPathPatterns := flag.Args()
	importPaths := gotool.ImportPaths(importPathPatterns)
	importPathsSet := make(map[string]bool, len(importPaths))
	for _, importPath := range importPaths {
		importPathsSet[importPath] = true
	}

	forward, _, graphErrors := BuildNoTests(&build.Default)
	if graphErrors != nil {
		log.Fatalln("importgraph.Build:", graphErrors)
	}

	renderGraph := func() ([]byte, error) {
		var in bytes.Buffer

		fmt.Fprintf(&in, "digraph \"\" {\n")
		for k, v := range forward {
			for k2 := range v {
				if !importPathsSet[k] || !importPathsSet[k2] || k == k2 {
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
		log.Fatalln("renderGraph:", err)
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
