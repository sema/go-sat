go-sat
======

Toy SAT solver written in the Go programming language.

go-sat implements the [DPLL algorithm](http://en.wikipedia.org/wiki/DPLL_algorithm) (or something resembling it).

To compile the ``go-sat`` executable run:

> go build

Usage (solving a SAT problem in CNF format):

> ./go-sat fixtures/small/sat/test.cnf
