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

package msg_test

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/internal/analyzer/msg"
	"fillmore-labs.com/zerolint/pkg/internal/diag"
)

func TestFormatf(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := [...]struct {
		name        string
		cat         diag.Category
		valueMethod bool
		format      string
		args        []any
		wantMessage string
	}{
		{
			name:        "simple message without args, no value method",
			cat:         CatArgumentNil,
			valueMethod: false,
			format:      "Hello",
			args:        []any{},
			wantMessage: "Hello (zl:arg)",
		},
		{
			name:        "simple message without args, with value method",
			cat:         CatArgumentNil,
			valueMethod: true,
			format:      "Hello",
			args:        []any{},
			wantMessage: "Hello (zl:arg+)",
		},
		{
			name:        "message with args, no value method",
			cat:         CatReceiver,
			valueMethod: false,
			format:      "Hello, %s",
			args:        []any{"World"},
			wantMessage: "Hello, World (zl:rcv)",
		},
		{
			name:        "message with args, with value method",
			cat:         CatReceiver,
			valueMethod: true,
			format:      "Number %d",
			args:        []any{42},
			wantMessage: "Number 42 (zl:rcv+)",
		},
		{
			name:        "empty format string, no value method",
			cat:         CatArgumentNil,
			valueMethod: false,
			format:      "",
			args:        []any{},
			wantMessage: " (zl:arg)",
		},
		{
			name:        "empty format string, with value method",
			cat:         CatArgumentNil,
			valueMethod: true,
			format:      "",
			args:        []any{},
			wantMessage: " (zl:arg+)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Formatf(tt.cat, tt.valueMethod, tt.format, tt.args...)
			if got.Message != tt.wantMessage {
				t.Errorf("Formatf() Message = %q, want %q", got.Message, tt.wantMessage)
			}

			if got.Category != tt.cat {
				t.Errorf("Formatf() Category = %q, want %q", got.Category, tt.cat)
			}
		})
	}
}

func TestFormatter(t *testing.T) {
	t.Parallel()

	mockType := types.Typ[types.String]
	ident1 := &ast.Ident{NamePos: token.NoPos, Name: "one"}
	ident2 := &ast.Ident{NamePos: token.NoPos, Name: "two"}

	tests := [...]struct {
		name      string
		formatter Formatter
	}{
		{
			name:      "func param",
			formatter: Param{},
		},
		{
			name:      "func result",
			formatter: Result{},
		},
		{
			name:      "struct field",
			formatter: Struct{},
		},
		{
			name:      "value spec",
			formatter: Value{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for _, names := range [][]*ast.Ident{
				nil,
				{ident1},
				{ident1, ident2},
			} {
				got := FormatMessage(tt.formatter, mockType, true, names)

				for _, name := range names {
					if !strings.Contains(got.Message, name.Name) {
						t.Errorf("FormatMessage() returned %q, does not contain %q", got.Message, name.Name)
					}
				}

				if len(got.Category) == 0 {
					t.Errorf("FormatMessage() Category is empty")
				}
			}
		})
	}
}
