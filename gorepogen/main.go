// gorepogen generates boilerplate files for Go repositories hosted on GitHub.
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

	"github.com/shurcooL/go/ioutil"
)

func t(text string) *template.Template {
	return template.Must(template.New("").Parse(text))
}

// Filename -> Template.
var templates = map[string]*template.Template{

	"README.md": t(`# {{.Title}} [![Build Status](https://travis-ci.org/{{.TravisCIPath}}.svg?branch=master)](https://travis-ci.org/{{.TravisCIPath}}) [![GoDoc](https://godoc.org/{{.ImportPath}}?status.svg)](https://godoc.org/{{.ImportPath}})

{{with .Doc}}{{.Doc}}{{else}}...
{{end}}
Installation
------------

` + "```bash" + `
go get -u {{.ImportPath}}{{if .NoGo}}/...{{end}}
{{if .HasJsTag}}go get -u -d -tags=js {{.ImportPath}}
{{end}}` + "```" + `
{{if not .HasLicenseFile}}
License
-------

-	[MIT License](https://opensource.org/licenses/mit-license.php)
{{end}}`),

	".travis.yml": t(`sudo: false
language: go
go:
  - 1.6.2
  - tip
matrix:
  allow_failures:
    - go: tip
  fast_finish: true
install:
  - # Do nothing. This is needed to prevent default install action "go get -t -v ./..." from happening here (we want it to happen inside script step).
script:
  - go get -t -v ./...
  - diff -u <(echo -n) <(gofmt -d -s .)
  - go tool vet .
  - go test -v -race ./...
`),
}

type goRepo struct {
	bpkg *build.Package
	NoGo bool
	Doc  *doc.Package
}

// ImportPath returns the import path for a GitHub repository.
func (r goRepo) ImportPath() string {
	return r.bpkg.ImportPath
}

// TravisCIPath returns the Travis CI path for a GitHub repository.
func (r goRepo) TravisCIPath() (string, error) {
	c := strings.Split(r.bpkg.ImportPath, "/")
	if len(c) < 3 {
		return "", errors.New("unexpected number of import path components")
	}
	if c[0] != "github.com" {
		return "", errors.New("Travis CI only supports GitHub repositories")
	}
	return path.Join(c[1], c[2]), nil
}

// Title is the package name for libraries and import path base for commands.
// TODO: And repo name otherwise.
func (r goRepo) Title() string {
	switch {
	case r.NoGo:
		return path.Base(r.bpkg.ImportPath)
	case r.bpkg.IsCommand():
		return path.Base(r.bpkg.ImportPath)
	case !r.bpkg.IsCommand():
		return r.bpkg.Name
	default:
		panic("unreachable")
	}
}

func (r goRepo) HasJsTag() bool {
	if r.NoGo { // TODO: Look in inner Go packages, if any?
		return false
	}
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
	var goRepo goRepo
	if bpkg, err := build.ImportDir(wd, build.ImportComment); err == nil {
		goRepo.bpkg = bpkg

		dpkg, err := docPackage(bpkg)
		if err != nil {
			return err
		}
		goRepo.Doc = dpkg
	} else if _, ok := err.(*build.NoGoError); ok {
		goRepo.bpkg = bpkg
		goRepo.NoGo = true
	} else {
		return err
	}

	for filename, t := range templates {
		var buf bytes.Buffer
		err = t.Execute(&buf, goRepo)
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

func main() {
	flag.Parse()

	err := gen()
	if err != nil {
		log.Fatalln(err)
	}
}
