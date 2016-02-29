// gopathshadow reports if you have any shadowed Go packages in your GOPATH workspaces.
//
// A shadowed package occurs when a Go package directory exists in more than one workspace,
// so one location will shadow others. This is an invalid setup and should be fixed.
package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	"github.com/kisielk/gotool"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: gopathshadow")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	clean := true

	workspaces := []string{build.Default.GOROOT}
	workspaces = append(workspaces, filepath.SplitList(build.Default.GOPATH)...)

	importPaths := gotool.ImportPaths([]string{"all"})
	for _, importPath := range importPaths {
		var dirs []string

		// Check which workspaces contain the directory for this Go package.
		for _, workspace := range workspaces {
			path := filepath.Join(workspace, "src", filepath.FromSlash(importPath))
			fi, err := os.Stat(path)
			if err == nil && fi.IsDir() {
				dirs = append(dirs, path)
			}
		}

		// If there's more than 1, then something is shadowed.
		if len(dirs) != 1 {
			clean = false

			fmt.Printf("warning: %q has %d dirs:\n", importPath, len(dirs))
			for _, path := range dirs {
				fmt.Printf("\t%s\n", path)
			}
		}
	}

	if !clean {
		os.Exit(1)
	}
}
