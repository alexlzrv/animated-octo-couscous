package main

import (
	"github.com/mayr0y/animated-octo-couscous.git/internal/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

// This package defines the main function for an analysis driver
// with several analyzers from packages go/analysis and staticcheck.io
// Usage:
// `go run cmd/staticlint/main.go <analyzers> <files>`

// Examples:
// Checking all files in the current folder with all analyzers
// go run cmd/staticlint/main.go ./...

func main() {
	multichecker.Main(
		GetAnalyzers()...,
	)
}

func GetAnalyzers() []*analysis.Analyzer {
	var analyzersSlice []*analysis.Analyzer
	analyzersSlice = append(analyzersSlice, analyzer.GetAnalysis()...)
	analyzersSlice = append(analyzersSlice, analyzer.GetStaticCheckAnalyzers()...)
	analyzersSlice = append(analyzersSlice, analyzer.ExitCheck)

	return analyzersSlice
}
