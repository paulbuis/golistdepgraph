# golistdepgraph

Rework of [godepgraph](https://github.com/kisielk/godepgraph) using execing `go list -json`
to gather info about packages.

## Install

    go get github.com/paulbuis/golistdepgraph

## Use

For basic usage, just give the package path of interest as the first
argument:

    golistdepgraph github.com/paulbuis/golistdepgraph

The output is a graph in [Graphviz](https://www.graphviz.org/) dot format. If you have the
graphviz tools installed you can render it by piping the output to `dot`:

    golistdepgraph github.com/paulbuis/golistdepgraph | dot -Tpng -o godepgraph.png

By default golistdepgraph will display packages in the standard library in the
graph, though it will not delve in to their dependencies.

### Ignoring Imports

#### The Go Standard Library

If you want to ignore standard library packages entirely, use the `-s` flag:

    golistdepgraph -s github.com/paulbuis/golistdepgraph

### 3 By Name

Import paths can be included in a comma-separated list passed to the `-i` flag:

    golistdepgraph -i github.com/foo/bar,github.com/baz/blah github.com/something/else

The packages and their imports will be excluded from the graph, unless the imports
are also imported by another package which is not excluded.

#### By Prefix

Import paths can also be ignored by prefix. The `-p` flag takes a comma-separated
list of prefixes:

    golistdepgraph -p github.com,launchpad.net bitbucket.org/foo/bar

### Delving into internals of GOROOT packages

    golistdepgraph -d github.com/paulbuis/golistdepgraph

### Including test packages

    golistdepgraph -t github.com/paulbuis/golistdepgraph

### Build tags

The current version of golistdepgraph simply execs the `go list -json` command and uses whatever it produces with the
current environment variables. It does not pass any build tags on to the `go` command line, although that is planned for
a future release.

## Output

### Node Background Fill Colors

golistdepgraph uses a simple color scheme to denote different types of packages:

  * *green*: a package that is part of the Go standard library, installed in `$GOROOT`.
  * *blue*: a regular Go package found in `$GOPATH`.
  * *orange*: a package found in `$GOPATH` that uses cgo by importing the special package "C".
  
  The above list of colors is incorrect!!!
  
### Node Label Font Colors

golistdepgraph uses a simple color scheme to denote different types of states packages, by font color:

  * *red*: an incomplete package that had an error in at least one dependency (e.g., missing source code).
  * *blue*: a stale package whose sources are not up-to-date with its binary.
  * *black*: no errors and up-to-date.

