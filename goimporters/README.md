goimporters
===========

Displays an import graph of Go packages that import the specified Go package in your GOPATH workspace.

Installation
------------

```bash
go get -u github.com/shurcooL/cmd/goimporters
```

Note that it requires `dot` command to be available (`brew install graphviz`).

Usage
-----

```
Usage: goimporters package
```

Example
-------

Here's a sample run:

```bash
goimporters "github.com/shurcooL/go/github_flavored_markdown"
```

Output:

![image](https://cloud.githubusercontent.com/assets/1924134/4436371/9442cd46-4774-11e4-9acb-500ac37c07a3.png)
