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
	"go/doc"
	"log"
	"os"
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

	"./README.md": t(`# {{.Doc.Name}} [![Build Status](https://travis-ci.org/{{.Username}}/{{.Doc.Name}}.svg?branch=master)](https://travis-ci.org/{{.Username}}/{{.Doc.Name}})

{{.Doc.Doc}}
Installation
------------

` + "```bash" + `
go get -u {{.Doc.ImportPath}}
` + "```" + `

License
-------

- [MIT License](http://opensource.org/licenses/mit-license.php)
`),

	"./.travis.yml": t(`language: go
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

func parseUsername(importPath string) (string, error) {
	c := strings.Split(importPath, "/")
	if len(c) < 3 {
		return "", errors.New("unexpected number of import path components")
	}
	return c[1], nil
}

func gen() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	dpkg, err := gist5504644.GetDocPackage(gist5504644.BuildPackageFromSrcDir(wd))
	if err != nil {
		return err
	}
	username, err := parseUsername(dpkg.ImportPath)
	if err != nil {
		return err
	}
	var state = struct {
		Doc      *doc.Package
		Username string
	}{
		Doc:      dpkg,
		Username: username,
	}

	for filename, t := range templates {
		var buf bytes.Buffer
		err = t.Execute(&buf, state)
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
