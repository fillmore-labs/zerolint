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
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/internal/passes/exclusions"
)

func TestAllDecl(t *testing.T) {
	t.Parallel()

	const src = `package main
func foo() {}
func bar() {}
`

	fs := token.NewFileSet()

	file, err := parser.ParseFile(fs, "main.go", src, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	var f *ast.FuncDecl
	for d := range AllDecl[*ast.FuncDecl]([]*ast.File{file}) {
		f = d

		break
	}

	if f == nil || f.Name.Name != "foo" {
		t.Errorf("Expected foo function, got %#v", f)
	}
}
