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

package exclusions

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// run performs the analysis to identify excluded types.
// It exports an [excludedFact] fact for each newly identified excluded type in the current package.
func run(pass *analysis.Pass) (any, error) {
	addExclusions(pass)

	for decl := range AllDecl[*ast.GenDecl](pass.Files) {
		if decl.Tok == token.TYPE {
			processTypeDecl(pass, decl)
		}
	}

	return newResult(pass), nil
}

func processTypeDecl(pass *analysis.Pass, genDecl *ast.GenDecl) {
	// Exclude all via "//zerolint:exclude" comment on the declaration block.
	excludeAll := HasExcludeComment(genDecl.Doc)

	for _, genSpec := range genDecl.Specs {
		spec, ok := genSpec.(*ast.TypeSpec)
		if !ok { // should not happen
			log.Printf("Internal error: Expected *ast.TypeSpec, got %T (zl:xxx)", genSpec)

			continue
		}

		if !excludeAll && !isCtype(spec) { // Exclude cgo-generated types.
			continue
		}

		def := pass.TypesInfo.Defs[spec.Name]
		tn, ok := def.(*types.TypeName)

		if !ok { // should not happen
			log.Printf("Internal error: Expected *types.TypeName, got %T (zl:xxx)", def)

			continue
		}

		excludeType(pass, tn)
	}
}

// isCtype is a heuristic to avoid issues with cgo types like _Ctype_struct_foo.
func isCtype(ts *ast.TypeSpec) bool {
	const cgoTypePrefix = "_Ctype_"

	return strings.HasPrefix(ts.Name.Name, cgoTypePrefix)
}

// excludeType exports an exclusion fact for the given object identifier.
func excludeType(pass *analysis.Pass, tn *types.TypeName) {
	pass.ExportObjectFact(tn, (*excludedFact)(nil))
}

// addExclusions prefills hard coded type definitions.
// For example, it ignores [runtime.Func] because pointers to this type represent opaque
// runtime-internal data, not zero-sized types the linter targets.
func addExclusions(pass *analysis.Pass) {
	pkg := pass.Pkg

	var typeNames []string

	// This might need an alternative approach:
	//
	// https://pkg.go.dev/golang.org/x/tools/go/analysis#hdr-Modular_analysis_with_Facts
	// “Some driver implementations (such as those based on Bazel and Blaze)
	// do not currently apply analyzers to packages of the standard library.”
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
			pass.ExportObjectFact(tn, (*excludedFact)(nil))
		}
	}
}
