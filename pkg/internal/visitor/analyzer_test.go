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

package visitor_test

import (
	"os"
	"slices"
	"testing"

	"fillmore-labs.com/zerolint/pkg/internal/excludes"
	"fillmore-labs.com/zerolint/pkg/internal/set"
	"fillmore-labs.com/zerolint/pkg/internal/visitor"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func TestRun(t *testing.T) {
	t.Parallel()

	dir := analysistest.TestData()

	ex, err := excludes.ReadExcludes(os.DirFS(dir), "excluded.txt")
	if err != nil {
		t.Fatalf("Can't find excludes file: %v", err)
	}

	type args struct {
		excludes set.Set[string]
		full     bool
		pattern  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"basic", args{nil, false, "go.test/basic"}, "go.test/basic.myError"},
		{"full", args{set.New(ex...), true, "go.test/a"}, "[0]string"},
	}

	for _, tt := range tests {
		v := visitor.New(visitor.Options{
			Excludes: tt.args.excludes,
			Full:     tt.args.full,
		})

		a := &analysis.Analyzer{
			Name:     "zerolint",
			Doc:      "...",
			Run:      v.Run,
			Requires: []*analysis.Analyzer{inspect.Analyzer},
		}
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			analysistest.RunWithSuggestedFixes(t, dir, a, tt.args.pattern)

			if !v.HasDetected() {
				t.Error("Expected detection of zero sized types")
			}

			zerotypes := slices.Collect(v.AllDetected())
			if !slices.Contains(zerotypes, tt.want) {
				t.Errorf("Expected %q to contain %q", zerotypes, tt.want)
			}
		})
	}
}
