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

package checker_test

import (
	"go/token"
	"go/types"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/internal/checker"
	"golang.org/x/tools/go/analysis"
)

type mockNode struct{}

func (mockNode) Pos() token.Pos {
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
		Pkg: types.NewPackage("example.com/test", "test"), // Important for calcIgnored
	}

	c := Checker{}
	c.Prepare(mockPass)

	rng := mockNode{}
	message := CategorizedMessage{
		Message:  "Test message (zl:test+)",
		Category: "test",
	}
	fixes := []analysis.SuggestedFix{
		{Message: "Fix1"},
	}

	c.Report(rng, message, fixes)

	if !called {
		t.Error("report was not called")
	}
}
