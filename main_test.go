package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func SatTester(t *testing.T, testdir string, shouldBeSat bool) {

	dir := path.Join(os.Getenv("GOPATH"), "src/github.com/sema/go-sat/fixtures", testdir)

	fileinfos, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Errorf("Could not open %v", dir)
		return
	}

	for _, fileInfo := range fileinfos {

		if fileInfo.IsDir() {
			continue
		}

		filePath := path.Join(dir, fileInfo.Name())
		_, sat, err := OpenAndSolve(filePath)

		if err != nil {
			t.Errorf("%v returned an error", filePath)
			continue
		}

		if sat != shouldBeSat {

			satStr := "sat"
			if shouldBeSat {
				satStr = "unsat"
			}

			t.Errorf("%v returned %v", filePath, satStr)
			continue
		}

		t.Logf("%v done", filePath)
	}

}

func TestSolverOnSmallSatExamples(t *testing.T) {
	SatTester(t, "small/sat", true)
}

func TestSolverOnSmallUnsatExamples(t *testing.T) {
	SatTester(t, "small/unsat", false)
}

func TestSolverOnHoos50SatExamples(t *testing.T) {
	SatTester(t, "hoos/50/sat", true)
}
