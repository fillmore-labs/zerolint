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

package excluded

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"fillmore-labs.com/zerolint/pkg/internal/filter"
	"fillmore-labs.com/zerolint/pkg/internal/set"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// run performs the analysis to identify excluded types.
// It collects types marked as excluded via:
// 1. Facts from dependencies ([excludeFact] facts).
// 2. Cgo-generated types (prefixed with _Ctype_).
// 3. Types with a "//zerolint:exclude" comment.
// It exports an [excludeFact] fact for each newly identified excluded type in the current package
// and returns a [filter.Filter] containing the token.Pos of all excluded type definitions.
func run(pass *analysis.Pass) (any, error) { //nolint:cyclop
	addExclusions(pass)

	in, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, ErrNoInspectorResult
	}

	for decl := range inspector.All[*ast.GenDecl](in) {
		if decl.Tok != token.TYPE {
			continue
		}

		// Exclude via "//zerolint:exclude" comment.
		excludeAll := hasExcludeComment(decl.Doc)

		for _, spec := range decl.Specs {
			if ts, ok := spec.(*ast.TypeSpec); ok {
				// Exclude cgo-generated types.
				if strings.HasPrefix(ts.Name.Name, "_Ctype_") {
					excludeType(pass, ts)

					continue
				}

				// Exclude via "//zerolint:exclude" comment.
				if excludeAll || hasExcludeComment(ts.Doc) || hasExcludeComment(ts.Comment) {
					excludeType(pass, ts)
				}
			}
		}
	}

	excludedTypeDefs := set.New[token.Pos]()

	for _, f := range pass.AllObjectFacts() {
		if _, ok := f.Fact.(*excludeFact); ok && f.Object != nil {
			excludedTypeDefs.Insert(f.Object.Pos())
		}
	}

	return filter.New(excludedTypeDefs), nil
}

// hasExcludeComment checks if a comment group contains "zerolint:exclude".
func hasExcludeComment(comments *ast.CommentGroup) bool {
	if comments == nil {
		return false
	}

	for _, comment := range comments.List {
		i := strings.Index(comment.Text, zerolintMarker)
		if i < 0 {
			continue
		}

		argsPart := strings.Fields(comment.Text[i+len(zerolintMarker):])
		if len(argsPart) == 0 {
			continue
		}

		args := strings.Split(argsPart[0], ",")
		if slices.Contains(args, "exclude") {
			return true
		}
	}

	return false
}

// excludeType exports an exclusion fact for the given object identifier.
func excludeType(pass *analysis.Pass, ts *ast.TypeSpec) {
	if tn := pass.TypesInfo.Defs[ts.Name]; tn != nil {
		pass.ExportObjectFact(tn, &excluded)
	}
}

// addExclusions prefills hard coded type definitions.
// For example, it ignores [runtime.Func] because pointers to this type represent opaque
// runtime-internal data, not zero-sized types the linter targets.
func addExclusions(pass *analysis.Pass) {
	pkg := pass.Pkg

	var typeNames []string

	switch pkg.Path() {
	case "runtime":
		typeNames = []string{"Func", "notInHeap"}

	case "runtime/cgo":
		typeNames = []string{"Incomplete"}

	default:
		return
	}

	scope := pkg.Scope()
	for _, name := range typeNames {
		if tn, ok := scope.Lookup(name).(*types.TypeName); ok {
			pass.ExportObjectFact(tn, &excluded)
		}
	}
}
