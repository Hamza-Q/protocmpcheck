package protocmpcheck

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func init() {
	debug = true
}

func TestFromFileSystem(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "testhelloworld")
}
