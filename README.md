cmd
===

[![Go Reference](https://pkg.go.dev/badge/github.com/shurcooL/cmd.svg)](https://pkg.go.dev/github.com/shurcooL/cmd)

Various small command-line utilities.

Installation
------------

```sh
go install github.com/shurcooL/cmd/...@latest
```

Directories
-----------

| Path                                                                                | Synopsis                                                                                                           |
|-------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------|
| [dumpargs](https://pkg.go.dev/github.com/shurcooL/cmd/dumpargs)                     | dumpargs dumps the command-line arguments.                                                                         |
| [dumpglfw3joysticks](https://pkg.go.dev/github.com/shurcooL/cmd/dumpglfw3joysticks) | dumpglfw3joysticks dumps state of attached joysticks.                                                              |
| [dumphttpreq](https://pkg.go.dev/github.com/shurcooL/cmd/dumphttpreq)               | dumphttpreq dumps incoming HTTP requests with full detail.                                                         |
| [godocrouter](https://pkg.go.dev/github.com/shurcooL/cmd/godocrouter)               | godocrouter is a reverse proxy that augments a private godoc server instance with global godoc.org instance.       |
| [goimporters](https://pkg.go.dev/github.com/shurcooL/cmd/goimporters)               | goimporters displays an import graph of Go packages that import the specified Go package in your GOPATH workspace. |
| [goimportgraph](https://pkg.go.dev/github.com/shurcooL/cmd/goimportgraph)           | goimportgraph displays an import graph within specified Go packages.                                               |
| [gopathshadow](https://pkg.go.dev/github.com/shurcooL/cmd/gopathshadow)             | gopathshadow reports if you have any shadowed Go packages in your GOPATH workspaces.                               |
| [gorepogen](https://pkg.go.dev/github.com/shurcooL/cmd/gorepogen)                   | gorepogen generates boilerplate files for Go repositories hosted on GitHub.                                        |
| [jsonfmt](https://pkg.go.dev/github.com/shurcooL/cmd/jsonfmt)                       | jsonfmt pretty-prints JSON from stdin.                                                                             |
| [runestats](https://pkg.go.dev/github.com/shurcooL/cmd/runestats)                   | runestats prints counts of total and unique runes from stdin.                                                      |
| [wasmserve](https://pkg.go.dev/github.com/shurcooL/cmd/wasmserve)                   | wasmserve compiles a Go command to WebAssembly and serves it via HTTP.                                             |

License
-------

-	[MIT License](LICENSE)
