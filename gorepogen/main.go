// Command gorepogen generates boilerplate files for Go repositories hosted on GitHub.
//
// Running it in repo root with a Go package writes files to the current working directory.
//
// It includes README.md with package doc, import path, MIT license, Travis badge,
// and .travis.yml that performs typical Go tests.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"go/doc"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/shurcooL/go/gists/gist5504644"
	"github.com/shurcooL/go/u/u11"
)

func t(text string) *template.Template {
	return template.Must(template.New("").Parse(text))
}

// Filename -> Template.
var templates = map[string]*template.Template{

	"README.md": t(`# {{.Title}} [![Build Status](https://travis-ci.org/{{.Username}}/{{.RepoName}}.svg?branch=master)](https://travis-ci.org/{{.Username}}/{{.RepoName}}) [![GoDoc](https://godoc.org/{{.Doc.ImportPath}}?status.svg)](https://godoc.org/{{.Doc.ImportPath}})

{{.Doc.Doc}}
Installation
------------

` + "```bash" + `
go get -u {{.Doc.ImportPath}}
{{if .HasJsTag}}go get -u -d -tags=js {{.Doc.ImportPath}}
{{end}}` + "```" + `
{{if not .HasLicenseFile}}
License
-------

- [MIT License](http://opensource.org/licenses/mit-license.php)
{{end}}`),

	".travis.yml": t(`language: go
go:
  - 1.4
install:
  - go get golang.org/x/tools/cmd/vet
script:
  - go get -t -v ./...
  - diff -u <(echo -n) <(gofmt -d ./)
  - go tool vet ./
  - go test -v -race ./...
`),
}

type goRepo struct {
	bpkg *build.Package
	Doc  *doc.Package
}

// Username is typically the 2nd import path element.
func (r goRepo) Username() (string, error) {
	c := strings.Split(r.Doc.ImportPath, "/")
	if len(c) < 3 {
		return "", errors.New("unexpected number of import path components")
	}
	return c[1], nil
}

// RepoName is the repository name, typically the 3rd import path element.
func (r goRepo) RepoName() (string, error) {
	c := strings.Split(r.Doc.ImportPath, "/")
	if len(c) < 3 {
		return "", errors.New("unexpected number of import path components")
	}
	return c[2], nil
}

// Title is the package name for libraries and import path base for commands.
func (r goRepo) Title() string {
	switch r.bpkg.IsCommand() {
	case true:
		return path.Base(r.bpkg.ImportPath)
	case false:
		return r.bpkg.Name
	}
	panic("unreachable")
}

func (r goRepo) HasJsTag() bool {
	for _, tag := range r.bpkg.AllTags {
		if tag == "js" {
			return true
		}
	}
	return false
}

// HasLicenseFile returns true if there's a LICENSE file present in current working directory.
func (r goRepo) HasLicenseFile() bool {
	if fi, err := os.Stat("LICENSE"); err == nil && !fi.IsDir() {
		return true
	}
	return false
}

func gen() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	bpkg, err := gist5504644.BuildPackageFromSrcDir(wd)
	if err != nil {
		return err
	}
	dpkg, err := gist5504644.GetDocPackage(bpkg, nil)
	if err != nil {
		return err
	}
	var goRepo = goRepo{
		bpkg: bpkg,
		Doc:  dpkg,
	}

	for filename, t := range templates {
		var buf bytes.Buffer
		err = t.Execute(&buf, goRepo)
		if err != nil {
			return err
		}
		fmt.Println("writing", filename)
		err = u11.WriteFile(&buf, filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	flag.Parse()

	err := gen()
	if err != nil {
		log.Fatalln(err)
	}
}
