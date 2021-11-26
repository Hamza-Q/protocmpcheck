package protocmpcheck

import "golang.org/x/tools/go/analysis"

var Analyzer = &analysis.Analyzer{
	Name:       "protocmpcheck",
	Doc:        "find calls to testify that compare protobuf types",
	Run:        run,
	ResultType: nil,
}

func run(pass *analysis.Pass) (interface{}, error) {
	return nil, nil
}
