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

package base_test

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/internal/base"
	"golang.org/x/tools/go/analysis"
)

// MockMsgFormatter is a mock implementation of the MsgFormatter interface for testing.
type MockMsgFormatter struct{}

func (m MockMsgFormatter) ZeroMsg(_ types.Type) (fmt.Stringer, string) {
	return mockCategory("zero"), "Zero message"
}

func (m MockMsgFormatter) SingularMsg(name string, _ types.Type) (fmt.Stringer, string) {
	return mockCategory("singular"), "Singular message for " + name
}

func (m MockMsgFormatter) PluralMsg(names string, _ types.Type) (fmt.Stringer, string) {
	return mockCategory("plural"), "Plural message for " + names
}

// mockCategory is a helper type to implement fmt.Stringer for testing.
type mockCategory string

func (m mockCategory) String() string {
	return string(m)
}

func TestFormatMessage(t *testing.T) {
	t.Parallel()

	mockFormatter := MockMsgFormatter{}
	mockType := types.Typ[types.Int]

	tests := []struct {
		name     string
		names    []*ast.Ident
		expected string
	}{
		{"ZeroNames", []*ast.Ident{}, "Zero message"},
		{"OneName", []*ast.Ident{{Name: "x"}}, "Singular message for x"},
		{"MultipleNames", []*ast.Ident{{Name: "x"}, {Name: "y"}}, `Plural message for "x", "y"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, msg := FormatMessage(tt.names, mockFormatter, mockType)
			if msg != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, msg)
			}
		})
	}
}

type mockNode struct{}

func (m mockNode) Pos() token.Pos {
	return token.Pos(1)
}

func (mockNode) End() token.Pos {
	return token.Pos(2)
}

func TestReport(t *testing.T) {
	t.Parallel()

	var called bool

	mockPass := &analysis.Pass{
		Report: func(diag analysis.Diagnostic) {
			t.Helper()

			called = true

			expectedMessage := "Test message (zl:test+)"
			if diag.Message != expectedMessage {
				t.Errorf("expected message %q, got %q", expectedMessage, diag.Message)
			}

			if diag.Category != "test" {
				t.Errorf("expected category %q, got %q", "test", diag.Category)
			}

			if len(diag.SuggestedFixes) != 1 || diag.SuggestedFixes[0].Message != "Fix1" {
				t.Errorf("unexpected suggested fixes: %+v", diag.SuggestedFixes)
			}
		},
		Pkg: types.NewPackage("test", "test"), // Important for calcIgnored
	}

	var base Base

	base.Prepare(mockPass)

	rng := mockNode{}
	category := mockCategory("test")
	message := "Test message"
	fixes := []analysis.SuggestedFix{
		{Message: "Fix1"},
	}

	base.Report(rng, category, true, message, fixes)

	if !called {
		t.Error("report was not called")
	}
}
