# cmd [![Build Status](https://travis-ci.org/shurcooL/cmd.svg?branch=master)](https://travis-ci.org/shurcooL/cmd) [![GoDoc](https://godoc.org/github.com/shurcooL/cmd?status.svg)](https://godoc.org/github.com/shurcooL/cmd)

Various small command-line utilities.

Installation
------------

```bash
go get -u github.com/shurcooL/cmd/...
```

Directories
-----------

| Path                                                                                   | Synopsis                                                                                                           |
|----------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------|
| [dump_args](https://godoc.org/github.com/shurcooL/cmd/dump_args)                       | dump_args dumps the command-line arguments.                                                                        |
| [dump_glfw3_joysticks](https://godoc.org/github.com/shurcooL/cmd/dump_glfw3_joysticks) | dump_glfw3_joysticks dumps state of attached joysticks.                                                            |
| [dump_httpreq](https://godoc.org/github.com/shurcooL/cmd/dump_httpreq)                 | dump_httpreq dumps incoming HTTP requests with full detail.                                                        |
| [godoc_router](https://godoc.org/github.com/shurcooL/cmd/godoc_router)                 | godoc_router is a reverse proxy that augments a private godoc server instance with global godoc.org instance.      |
| [goimporters](https://godoc.org/github.com/shurcooL/cmd/goimporters)                   | goimporters displays an import graph of Go packages that import the specified Go package in your GOPATH workspace. |
| [goimportgraph](https://godoc.org/github.com/shurcooL/cmd/goimportgraph)               | goimportgraph displays an import graph within specified Go packages.                                               |
| [gopathshadow](https://godoc.org/github.com/shurcooL/cmd/gopathshadow)                 | gopathshadow reports if you have any shadowed Go packages in your GOPATH workspaces.                               |
| [gorepogen](https://godoc.org/github.com/shurcooL/cmd/gorepogen)                       | gorepogen generates boilerplate files for Go repositories hosted on GitHub.                                        |
| [jsonfmt](https://godoc.org/github.com/shurcooL/cmd/jsonfmt)                           | jsonfmt pretty-prints JSON from stdin.                                                                             |
| [rune_stats](https://godoc.org/github.com/shurcooL/cmd/rune_stats)                     | rune_stats prints counts of total and unique runes from stdin.                                                     |
| [table](https://godoc.org/github.com/shurcooL/cmd/table)                               | table is a chef client command-line tool.                                                                          |

License
-------

-	[MIT License](https://opensource.org/licenses/mit-license.php)
