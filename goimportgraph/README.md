goimportgraph
=============

Displays an import graph within specified Go packages.

Installation
------------

```bash
go get -u github.com/shurcooL/cmd/goimportgraph
```

Note that it requires `dot` command to be available (`brew install graphviz`).

Usage
-----

```bash
goimportgraph packages
```

Example
-------

Here's a sample run:

```bash
goimportgraph encoding/...
```

Output:

![Screenshot](https://cloud.githubusercontent.com/assets/1924134/8923141/c6e54254-34a1-11e5-876c-4a8ab69feccb.png)
