package main

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
)

type pkg struct {
	ImportPath string
	Doc        string
}

// packagesInside returns a list of packages that have root as import path prefix,
// not including package with import path equal to root.
func packagesInside(root string) ([]pkg, error) {
	var pkgs []pkg
	cmd := exec.Command("go", "list", "-f", "{{.ImportPath}}\t{{.Doc}}", root+"/...")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	br := bufio.NewReader(stdout)
	for line, err := br.ReadString('\n'); err == nil; line, err = br.ReadString('\n') {
		importPathDoc := strings.Split(line[:len(line)-1], "\t") // Trim last newline.
		if len(importPathDoc) != 2 {
			log.Fatalf("len(importPathDoc) should be 2, but was %v; line was %q", len(importPathDoc), line)
		}
		importPath, doc := importPathDoc[0], importPathDoc[1]
		if importPath == root {
			continue
		}
		pkgs = append(pkgs, pkg{ImportPath: importPath, Doc: doc})
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return pkgs, nil
}
