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
		goRepo.Doc, err = docPackage(bpkg)
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

// TODO: Keep in sync or unify with github.com/shurcooL/vfsgen/cmd/vfsgendev/parse.go.
// TODO: See if these can be cleaned up.

func docPackage(bpkg *build.Package) (*doc.Package, error) {
	apkg, err := astPackage(bpkg)
	if err != nil {
		return nil, err
	}
	return doc.New(apkg, bpkg.ImportPath, 0), nil
}

func astPackage(bpkg *build.Package) (*ast.Package, error) {
	// TODO: Either find a way to use golang.org/x/tools/importer (from Go 1.4~ or older, it no longer exists as of Go 1.6) directly, or do file AST parsing in parallel like it does.
	filenames := append(bpkg.GoFiles, bpkg.CgoFiles...)
	files := make(map[string]*ast.File, len(filenames))
	fset := token.NewFileSet()
	for _, filename := range filenames {
		fileAst, err := parser.ParseFile(fset, filepath.Join(bpkg.Dir, filename), nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files[filename] = fileAst // TODO: Figure out if filename or full path are to be used (the key of this map doesn't seem to be used anywhere).
	}
	return &ast.Package{Name: bpkg.Name, Files: files}, nil
}
