package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"code.google.com/p/go.tools/refactor/importgraph"

	"github.com/shurcooL/go-goon"
	"github.com/shurcooL/go/u/u3"
)

func main() {
	forward, reverse, graphErrors := importgraph.Build(&build.Default)
	_, _, _ = forward, reverse, graphErrors
	if graphErrors != nil {
		goon.DumpExpr(graphErrors)
		panic(0)
	}
	//goon.DumpExpr(forward)
	//goon.DumpExpr(forward.Search("github.com/shurcooL/go/github_flavored_markdown"))

	target := "github.com/shurcooL/go/github_flavored_markdown"
	if len(os.Args) >= 2 {
		target = os.Args[1]
	}
	reachable := reverse.Search(target)

	renderGraph := func() ([]byte, error) {
		var in bytes.Buffer

		fmt.Fprintf(&in, "digraph \"\" {\n")
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
