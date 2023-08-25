package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/analysis/code"
)

var ExitCheck = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "checks for os.Exit call in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename

		if !strings.HasSuffix(filename, ".go") {
			continue
		}

		if strings.HasSuffix(filename, "_test.go") {
			continue
		}

		if code.IsMain(pass) {
			ast.Inspect(file, func(node ast.Node) bool {
				if _, ok := node.(*ast.CallExpr); ok && code.IsCallTo(pass, node, "os.Exit") && !code.IsInTest(pass, node) {
					pass.Reportf(node.Pos(), "call to os.Exit() in main")
				}
				return true
			})
		}
	}
	return nil, nil
}
