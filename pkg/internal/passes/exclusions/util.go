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
	"go/types"
	"iter"

	"golang.org/x/tools/go/analysis"
)

// AllDecls iterates over all declarations of the specified type T in the given files.
func AllDecls[T ast.Decl](files []*ast.File) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, f := range files {
			for _, decl := range f.Decls {
				if d, ok := decl.(T); ok {
					if !yield(d) {
						return
					}
				}
			}
		}
	}
}

// AllFacts iterates over all facts of the specified type T in the given facts.
func AllFacts[T analysis.Fact](facts []analysis.ObjectFact) iter.Seq2[types.Object, T] {
	return func(yield func(types.Object, T) bool) {
		for _, fact := range facts {
			if f, ok := fact.Fact.(T); ok {
				if !yield(fact.Object, f) {
					break
				}
			}
		}
	}
}
