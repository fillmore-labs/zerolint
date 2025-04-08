// Copyright 2024 Oliver Eikemeier. All Rights Reserved.
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

package excluded_test

import (
	"errors"
	"go/ast"
	"go/types"
	"testing"

	"fillmore-labs.com/zerolint/pkg/internal/filter"
	. "fillmore-labs.com/zerolint/pkg/internal/passes/excluded"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func TestExcludedAnalyzer(t *testing.T) {
	t.Parallel()

	dir := analysistest.TestData()
	analysistest.Run(t, dir, TestAnalyzer, "./...")
}

var TestAnalyzer = &analysis.Analyzer{ //nolint:gochecknoglobals
	Name:     "testanalyzer",
	Doc:      "consumes results from excluded.Analyzer for testing",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer, Analyzer},
}

var ErrNoExcludedResult = errors.New("result of excluded.Analyzer missing")

func run(pass *analysis.Pass) (any, error) {
	in, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, ErrNoInspectorResult
	}

	excludedResult, ok := pass.ResultOf[Analyzer].(filter.Filter)
	if !ok {
		return nil, ErrNoExcludedResult
	}

	for valueSpec := range inspector.All[*ast.ValueSpec](in) {
		typ := valueSpec.Type
		if named, ok := pass.TypesInfo.TypeOf(typ).(*types.Named); ok {
			tn := named.Obj()

			msg := "IS NOT excluded"
			if excludedResult.ExcludedType(tn) {
				msg = "IS excluded"
			}

			pass.Reportf(typ.Pos(), "Type %q %s", tn.Name(), msg)
		}
	}

	return any(nil), nil
}
