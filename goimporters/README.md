goimporters
===========

Displays an import graph of Go packages that import the specified Go package in your GOPATH workspace.

Installation
------------

```sh
go install github.com/shurcooL/cmd/goimporters@latest
```

Note that it requires `dot` command to be available (`brew install graphviz`).

Usage
-----

```sh
goimporters packages
```

Example
-------

Here's a sample run:

```sh
goimporters github.com/shurcooL/github_flavored_markdown
```

Output:

![Screenshot](https://cloud.githubusercontent.com/assets/1924134/4436371/9442cd46-4774-11e4-9acb-500ac37c07a3.png)
