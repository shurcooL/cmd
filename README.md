cmd
===

[![Build Status](https://travis-ci.org/shurcooL/cmd.svg?branch=master)](https://travis-ci.org/shurcooL/cmd) [![GoDoc](https://godoc.org/github.com/shurcooL/cmd?status.svg)](https://godoc.org/github.com/shurcooL/cmd)

Various small command-line utilities.

Installation
------------

```bash
go get -u github.com/shurcooL/cmd/...
```

Directories
-----------

| Path                                                                               | Synopsis                                                                                                           |
|------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------|
| [dumpargs](https://godoc.org/github.com/shurcooL/cmd/dumpargs)                     | dumpargs dumps the command-line arguments.                                                                         |
| [dumpglfw3joysticks](https://godoc.org/github.com/shurcooL/cmd/dumpglfw3joysticks) | dumpglfw3joysticks dumps state of attached joysticks.                                                              |
| [dumphttpreq](https://godoc.org/github.com/shurcooL/cmd/dumphttpreq)               | dumphttpreq dumps incoming HTTP requests with full detail.                                                         |
| [godocrouter](https://godoc.org/github.com/shurcooL/cmd/godocrouter)               | godocrouter is a reverse proxy that augments a private godoc server instance with global godoc.org instance.       |
| [goimporters](https://godoc.org/github.com/shurcooL/cmd/goimporters)               | goimporters displays an import graph of Go packages that import the specified Go package in your GOPATH workspace. |
| [goimportgraph](https://godoc.org/github.com/shurcooL/cmd/goimportgraph)           | goimportgraph displays an import graph within specified Go packages.                                               |
| [gopathshadow](https://godoc.org/github.com/shurcooL/cmd/gopathshadow)             | gopathshadow reports if you have any shadowed Go packages in your GOPATH workspaces.                               |
| [gorepogen](https://godoc.org/github.com/shurcooL/cmd/gorepogen)                   | gorepogen generates boilerplate files for Go repositories hosted on GitHub.                                        |
| [jsonfmt](https://godoc.org/github.com/shurcooL/cmd/jsonfmt)                       | jsonfmt pretty-prints JSON from stdin.                                                                             |
| [runestats](https://godoc.org/github.com/shurcooL/cmd/runestats)                   | runestats prints counts of total and unique runes from stdin.                                                      |
| [wasmserve](https://godoc.org/github.com/shurcooL/cmd/wasmserve)                   | wasmserve compiles a Go command to WebAssembly and serves it via HTTP.                                             |

License
-------

-	[MIT License](https://opensource.org/licenses/mit-license.php)
