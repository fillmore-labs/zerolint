// Copyright 2024-2025 Oliver Eikemeier. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package exclusions_test

import (
	"errors"
	"go/ast"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	. "fillmore-labs.com/zerolint/pkg/internal/passes/exclusions"
)

func TestExclusionsAnalyzer(t *testing.T) {
	t.Parallel()

	testAnalyzer := &analysis.Analyzer{
		Name:     "testanalyzer",
		Doc:      "consumes results from exclusions.Analyzer for testing",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer, Analyzer},
	}

	missingAnalyzer := &analysis.Analyzer{
		Name:     "missinganalyzer",
		Doc:      "missing results from exclusions.Analyzer",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	dir := analysistest.TestData()

	tests := []struct {
		name     string
		analyzer *analysis.Analyzer
		pkg      string
	}{
		{"exclusions", testAnalyzer, "test/a"},
		{"calculated", testAnalyzer, "test/e"},
		{"missing", missingAnalyzer, "test/n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			analysistest.Run(t, dir, tt.analyzer, tt.pkg)
		})
	}
}

var ErrNoInspectorResult = errors.New("testanalyzer: inspector result missing")

func run(pass *analysis.Pass) (any, error) {
	in, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, ErrNoInspectorResult
	}

	excludedTypeDefs, err := CalculateExclusions(pass)
	if err != nil && !errors.Is(err, ErrNoExclusionsResult) {
		return nil, err
	}

	for valueSpec := range inspector.All[*ast.ValueSpec](in) {
		t := pass.TypesInfo.TypeOf(valueSpec.Type)

		var tn *types.TypeName

		switch t := t.(type) {
		case *types.Named:
			tn = t.Obj()

		case *types.Alias:
			tn = t.Obj()

		default:
			continue
		}

		msg := "IS NOT excluded"
		if excludedTypeDefs.Has(tn.Pos()) {
			msg = "IS excluded"
		}

		pass.Reportf(valueSpec.Pos(), "Type %q %s", tn.Name(), msg)
	}

	return any(nil), nil
}
