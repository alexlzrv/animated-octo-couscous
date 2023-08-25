package analyzer_test

import (
	"testing"

	"github.com/mayr0y/animated-octo-couscous.git/internal/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestExitAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.ExitCheck, "./...")
}
