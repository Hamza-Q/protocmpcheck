package main

import (
	"github.com/Hamza-Q/protocmpcheck"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(protocmpcheck.Analyzer)
}
