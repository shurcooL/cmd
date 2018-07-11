// gorepogen generates boilerplate files for Go repositories hosted on GitHub.
//
// Running it in repo root with a Go package writes files to the current working directory.
//
// It includes README.md with package doc, import path, MIT license, Travis badge,
// and .travis.yml that performs typical Go tests.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/shurcooL/go/ioutil"
)

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	var goRepo goRepo
	if bpkg, err := build.ImportDir(wd, build.ImportComment); err == nil {
		goRepo.bpkg = bpkg
		goRepo.Doc, err = computeDoc(bpkg)
		if err != nil {
			return err
		}
	} else if _, ok := err.(*build.NoGoError); ok {
		goRepo.bpkg = bpkg
		goRepo.NoGo = true
	} else {
		return err
	}

	for filename, t := range templates {
		var buf bytes.Buffer
		err := t.Execute(&buf, goRepo)
		if err != nil {
			return err
		}
		fmt.Println("writing", filename)
		err = ioutil.WriteFile(filename, &buf)
		if err != nil {
			return err
		}
	}

	return nil
}

// computeDoc computes the package documentation for the given package.
func computeDoc(bpkg *build.Package) (*doc.Package, error) {
	fset := token.NewFileSet()
	files := make(map[string]*ast.File)
	for _, file := range append(bpkg.GoFiles, bpkg.CgoFiles...) {
		f, err := parser.ParseFile(fset, filepath.Join(bpkg.Dir, file), nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files[file] = f
	}
	apkg := &ast.Package{
		Name:  bpkg.Name,
		Files: files,
	}
	return doc.New(apkg, bpkg.ImportPath, 0), nil
}
