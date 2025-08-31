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

package analyzer_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/ast/inspector"

	. "fillmore-labs.com/zerolint/pkg/internal/analyzer"
	"fillmore-labs.com/zerolint/pkg/internal/typeutil"
)

func TestAllReturns(t *testing.T) { //nolint:funlen
	t.Parallel()

	const src = `package main

func oneReturn() int {
	return 1
}

func twoReturns(x int) int {
	if x > 0 {
		return 1
	}
	return 0
}

func noReturns() {
	// no-op
}

func nestedReturn() {
	_ = func() int {
		return 42 // This should be ignored
	}
	return // This should be found
}
`
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "main.go", src, 0)
	if err != nil {
		t.Fatalf("failed to parse source: %v", err)
	}

	files := []*ast.File{file}
	root := inspector.New(files).Root()

	testCases := []struct {
		funcName string
		want     int
	}{
		{"oneReturn", 1},
		{"twoReturns", 2},
		{"noReturns", 0},
		{"nestedReturn", 1},
	}

	for _, tc := range testCases {
		t.Run(tc.funcName, func(t *testing.T) {
			t.Parallel()

			// Find the function declaration we want to test
			var targetNode ast.Node

			for fn := range typeutil.AllDecls[*ast.FuncDecl](files) {
				if fn.Name.Name == tc.funcName {
					targetNode = fn

					break
				}
			}

			if targetNode == nil {
				t.Fatalf("function %q not found in test source", tc.funcName)
			}

			c, ok := root.FindNode(targetNode)
			if !ok {
				t.Fatalf("failed to find cursor for function %q", tc.funcName)
			}

			count := 0
			for range AllReturns(c) {
				count++
			}

			if count != tc.want {
				t.Errorf("AllReturns() found %d returns, want %d", count, tc.want)
			}
		})
	}
}
