package main

import (
	"github.com/mayr0y/animated-octo-couscous.git/internal/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

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
