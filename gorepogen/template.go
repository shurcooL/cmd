package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"go/doc"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/mattn/go-runewidth"
	"github.com/shurcooL/markdownfmt/markdown"
)

// Filename -> Template.
var templates = map[string]*template.Template{

	"README.md": t(`{{.Title | underline}}

[![Build Status](https://travis-ci.org/{{.TravisCIPath}}.svg?branch=master)](https://travis-ci.org/{{.TravisCIPath}}) [![GoDoc](https://godoc.org/{{.ImportPath}}?status.svg)](https://godoc.org/{{.ImportPath}})

{{with .Doc}}{{.Doc}}{{else}}...
{{end}}
Installation
------------

` + "```bash" + `
go get -u {{.ImportPath}}{{if .NoGo}}/...{{end}}
{{if .HasJSTag}}GOARCH=js go get -u -d {{.ImportPath}}
{{end}}` + "```" + `
{{with .Directories}}
Directories
-----------

{{.}}{{end}}
License
-------
{{if .HasLicenseFile}}
-	[MIT License](LICENSE)
{{else}}
-	[MIT License](https://opensource.org/licenses/mit-license.php)
{{end}}`),

	".travis.yml": t(`sudo: false
language: go
go:
  - 1.x
  - master
matrix:
  allow_failures:
    - go: master
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

func (r goRepo) HasJSTag() bool {
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

func (r goRepo) Directories() (string, error) {
	pkgs, err := packagesInside(r.bpkg.ImportPath)
	if err != nil {
		return "", err
	}

	// If there are no packages, don't include a directories section.
	if len(pkgs) == 0 {
		return "", nil
	}

	md := new(bytes.Buffer)
	fmt.Fprintln(md, "Path | Synopsis")
	fmt.Fprintln(md, "-----|---------")
	for _, p := range pkgs {
		relativePath := strings.TrimPrefix(p.ImportPath, r.bpkg.ImportPath+"/")
		fmt.Fprintf(md, "[%s](%s) | %s\n", relativePath, "https://godoc.org/"+p.ImportPath, p.Doc)
	}

	formatted, err := markdown.Process("", md.Bytes(), nil)
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

type pkg struct {
	ImportPath string
	Doc        string
}

// packagesInside returns a list of packages that have root as import path prefix,
// not including package with import path equal to root.
func packagesInside(root string) ([]pkg, error) {
	cmd := exec.Command("go", "list", "-f", "{{.ImportPath}}\t{{.Doc}}", root+"/...")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	var pkgs []pkg
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
	err = cmd.Wait()
	if err != nil {
		return nil, err
	}
	return pkgs, nil
}

// HasLicenseFile returns true if there's a LICENSE file present in current working directory.
func (r goRepo) HasLicenseFile() bool {
	fi, err := os.Stat("LICENSE")
	return err == nil && !fi.IsDir()
}

func t(text string) *template.Template {
	return template.Must(template.New("").Funcs(template.FuncMap{
		"underline": func(s string) string { return s + "\n" + strings.Repeat("=", runewidth.StringWidth(s)) },
	}).Parse(text))
}
