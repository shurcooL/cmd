// goimportgraph displays an import graph within specified Go packages.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/kisielk/gotool"
	"github.com/shurcooL/go/importgraphutil"
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

var (
	testsFlag = flag.Bool("tests", false, "Include tests when building graph.")
)

func usage() {
	fmt.Fprint(os.Stderr, "Usage: goimportgraph [packages]\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, `
Examples:
  goimportgraph encoding/...

  goimportgraph github.com/shurcooL/Conception-go/...
`)
}

func init() { log.SetFlags(0) }

func main() {
	flag.Usage = usage
	flag.Parse()

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	importPathPatterns := flag.Args()
	if len(importPathPatterns) == 0 {
		importPathPatterns = []string{"."}
	}
	importPaths := gotool.ImportPaths(importPathPatterns)
	importPaths, err := resolveRelative(importPaths)
	if err != nil {
		return err
	}
	importPathsSet := make(map[string]bool, len(importPaths))
	for _, importPath := range importPaths {
		importPathsSet[importPath] = true
	}

	var forward importgraph.Graph
	var graphErrors map[string]error
	switch *testsFlag {
	case false:
		forward, _, graphErrors = importgraphutil.BuildNoTests(&build.Default)
	case true:
		forward, _, graphErrors = importgraph.Build(&build.Default)
	}
	if graphErrors != nil {
		return fmt.Errorf("importgraph.Build: %v", graphErrors)
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
		return cmd.Output()
	}

	graphSvg, err := renderGraph()
	if err != nil {
		return fmt.Errorf("renderGraph: %v", err)
	}

	mux := http.NewServeMux()
	stopServerChan := make(chan struct{})
	mux.HandleFunc("/index.svg", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(graphSvg)

		go func() {
			time.Sleep(time.Second)
			stopServerChan <- struct{}{}
		}()
	})
	u3.DisplayHtmlInBrowser(mux, stopServerChan, ".svg")

	return nil
}

// resolveRelative checks that the packages exist, resolves relative import paths to full.
func resolveRelative(importPaths []string) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for i, path := range importPaths {
		bpkg, err := build.Import(path, wd, 0)
		if err != nil {
			return nil, fmt.Errorf("can't load package %q: %v", path, err)
		}
		importPaths[i] = bpkg.ImportPath
	}
	return importPaths, nil
}
