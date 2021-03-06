package protocmpcheck

import (
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"

	// These proto repos do not need to be directly imported, but are imported
	// by name via the packages importer.
	// Include them as imports so that they are identified and versioned by
	// go.mod.
	_ "github.com/golang/protobuf/proto"
	_ "google.golang.org/protobuf/proto"
)

var debug = false

func logf(msg string, args ...interface{}) {
	if debug {
		log.Printf(msg, args...)
	}
}

var Analyzer = &analysis.Analyzer{
	Name: "protocmpcheck",
	Doc:  "find calls to testify that compare protobuf types",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				checkCall(pass, call)
			}
			return true
		})
	}
	return nil, nil
}

// Checks the call for invalid protobuf comparisons.
// Returns false if reported an error.
func checkCall(pass *analysis.Pass, call *ast.CallExpr) (shouldRecurse bool) {
	shouldRecurse = true
	// Ignore calls to locally defined functions - they are almost definitely
	// not calls to testify.
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	// Type conversions cannot be calls to testify functions.
	if pass.TypesInfo.Types[sel.Sel].IsType() {
		return
	}
	logf("Found SelectorExpr: %s", sel.Sel.Name)

	selIdent, _ := sel.X.(*ast.Ident)
	logf("with selIdent: %s", selIdent)
	selObj, ok := pass.TypesInfo.Uses[selIdent]
	if !ok {
		logf("no uses")
		return
	}
	pkgName, ok := selObj.(*types.PkgName)
	var fullImportPath string
	if ok {
		fullImportPath = pkgName.Imported().Path()
	} else {
		logf("not a pkgName: %s", selIdent)
		// Extract package if this is a method.
		// We know testify uses pointer receivers so ignore value receivers to
		// avoid complexity
		ptr, ok := selObj.Type().(*types.Pointer)
		if !ok {
			logf("not pointer")
			return
		}
		ptrType, ok := ptr.Elem().(*types.Named)
		if !ok {
			logf("not named")
			return
		}
		fullImportPath = ptrType.Obj().Pkg().Path()
		logf(fullImportPath)
	}
	if !isComparisonPackage(fullImportPath) {
		logf("not comparison pkg")
		return
	}
	return checkComparison(pass, call, sel, fullImportPath)
}

// is the selector expression calling a package known to make comparisons?
func isComparisonPackage(pkgName string) bool {
	switch pkgName {
	case "reflect", "github.com/stretchr/testify/require", "github.com/stretchr/testify/assert":
		return true
	default:
		return false
	}
}

func checkComparison(pass *analysis.Pass, call *ast.CallExpr, sel *ast.SelectorExpr, pkgName string) (shouldRecurse bool) {
	shouldRecurse = true
	if pkgName == "reflect" && sel.Sel.Name != "DeepEqual" {
		return
	} else if strings.Contains(pkgName, "testify") && !isCheckableTestifyCall(sel.Sel.Name) {
		return
	}

	var arg1, arg2 ast.Expr
	if pkgName == "reflect" || isMethod(pass, sel) {
		arg1, arg2 = call.Args[0], call.Args[1]
	} else {
		arg1, arg2 = call.Args[1], call.Args[2]
	}
	if checkArgs(pass, arg1, arg2) {
		pass.Report(analysis.Diagnostic{
			Pos:     call.Lparen,
			Message: fmt.Sprintf("call to %s.%s comparing protobuf types not allowed", pkgName, sel.Sel.Name),
		})
		return false
	}
	return true
}

func isMethod(pass *analysis.Pass, sel *ast.SelectorExpr) (ok bool) {
	defer func() {
		logf("isMethod: %+v %+v {%v}", sel.X, sel.Sel, ok)
	}()
	t, ok := pass.TypesInfo.Types[sel.X]
	if !ok {
		logf("not method")
		return false
	}
	return strings.Contains(t.Type.String(), "Assertion")
}

func isCheckableTestifyCall(funcName string) bool {
	checkable := []string{
		"Equal",
		"ElementsMatch",
		"EqualValues",
	}
	for _, c := range checkable {
		if funcName == c {
			return true
		}
	}
	return false
}

// return true if args are problematic
func checkArgs(pass *analysis.Pass, arg1, arg2 ast.Expr) bool {
	return isProto(pass, arg1) && isProto(pass, arg2)
}

var protoMessageTypes []*types.Interface

// TODO: Move to top of file
func init() {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedImports,
	}
	pkgs, err := packages.Load(cfg, "google.golang.org/protobuf/proto", "github.com/golang/protobuf/proto")
	if err != nil {
		panic(err)
	}
	if len(pkgs) != 2 {
		panic("lookup fialed")
	}
	for _, pkg := range pkgs {
		protoType := pkg.Types.Scope().Lookup("Message").Type().Underlying().(*types.Interface)
		protoMessageTypes = append(protoMessageTypes, protoType)
	}
}

func isProto(pass *analysis.Pass, arg ast.Expr) bool {
	argType := pass.TypesInfo.Types[arg].Type
	logf("arg: %s type: %s", arg, argType)

	for _, pType := range protoMessageTypes {
		if types.Implements(argType, pType) {
			return true
		}
	}
	return false
}
