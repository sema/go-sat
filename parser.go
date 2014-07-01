package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseBenchmarkFile(path string) (formula Formula, err error) {

	// Open file

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	// Scan file

	var insertedClauses = 0

	var scanner = bufio.NewScanner(file)
	for scanner.Scan() {

		var line = scanner.Text()
		var ts = strings.Split(line, " ")
		var tokens = make([]string, len(ts))

		// some benchmark sets inject random spaces.
		// filter out all "empty" tokens

		var top = 0
		for _, t := range ts {
			if t != "" {
				tokens[top] = t
				top++
			}
		}

		tokens = tokens[:top]

		if len(tokens) == 0 || tokens[0] == "c" || tokens[0] == "%" || tokens[0] == "0" {
			continue // ignore blank lines or comments

		} else if len(tokens) == 4 && tokens[0] == "p" && tokens[1] == "cnf" {

			var _, err1 = strconv.Atoi(tokens[2])
			var numClauses, err2 = strconv.Atoi(tokens[3])

			if err1 == nil && err2 == nil {
				formula = make(Formula, numClauses)
				continue
			}

		} else if formula != nil {

			var clause = make(Clause, len(tokens)-1)

			for i, token := range tokens {

				var atom int

				atom, err = strconv.Atoi(token)
				if err != nil {
					return
				}

				var isNegated = false
				if atom < 0 {
					isNegated = true
					atom = atom * -1
				}

				if atom != 0 {
					var literal Literal
					literal.atom = Atom(atom)
					literal.negated = isNegated

					clause[i] = literal
				}

			}

			formula[insertedClauses] = clause
			insertedClauses++

			continue
		}

		err = errors.New(fmt.Sprint("Parsing error at line: ", line, "\n"))
		return

	}

	err = scanner.Err()
	return

}
