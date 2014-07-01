
package main

import (
    "fmt"
    "os"
)

func OpenAndSolve(path string) (solution Solution, sat bool, err error) {
 
    formula, err := ParseBenchmarkFile(path)
    if err != nil {
        return
    }

    solution, sat = Solve(formula)
    return
}

func main() {
    if len(os.Args) != 2 {
        fmt.Printf("Usage: %v <cnf formatted file>\n", os.Args[0])
        os.Exit(1)
    }

    solution, _, err := OpenAndSolve(os.Args[1])

    if err != nil {
        fmt.Printf("Error: %v", err.Error())
        os.Exit(1)
    }

    fmt.Print(solution)
}
