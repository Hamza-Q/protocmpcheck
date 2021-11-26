package protocmpcheck

import (
	"testing"

	_ "github.com/Hamza-Q/protocmpcheck/testdata/src/testhelloworld"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFromFileSystem(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "testhelloworld")
}
