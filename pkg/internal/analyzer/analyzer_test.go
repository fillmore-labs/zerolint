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
	"os"
	"reflect"
	"regexp"
	"slices"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"

	. "fillmore-labs.com/zerolint/pkg/internal/analyzer"
	"fillmore-labs.com/zerolint/pkg/internal/excludes"
	"fillmore-labs.com/zerolint/pkg/internal/passes/exclusions"
	"fillmore-labs.com/zerolint/pkg/internal/set"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
	"fillmore-labs.com/zerolint/pkg/zerolint/result"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	dir := analysistest.TestData()

	excludedTypeNames, err := excludes.ReadExcludes(os.DirFS(dir), "excluded.txt")
	if err != nil {
		t.Fatalf("Can't find excludes file: %v", err)
	}

	type args struct {
		options Options
		pkg     string
	}

	testre := regexp.MustCompile("^test/.*$")

	tests := []struct {
		name string
		args args
		want string
	}{
		{"basic", args{Options{Level: level.Basic, Regex: testre}, "test/basic"}, "test/basic.myError (value methods)"},
		{"full", args{Options{Level: level.Full, Excludes: set.New(excludedTypeNames...)}, "test/a"}, "[0]string"},
		{"exclusions", args{Options{Level: level.Full}, "test/e"}, ""},
	}
	for _, tt := range tests {
		a := &analysis.Analyzer{
			Name:       "zerolint",
			Doc:        "...",
			Run:        tt.args.options.Run,
			Requires:   []*analysis.Analyzer{inspect.Analyzer, exclusions.Analyzer},
			ResultType: reflect.TypeFor[result.Detected](),
		}

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := analysistest.RunWithSuggestedFixes(t, dir, a, tt.args.pkg)

			if tt.want != "" {
				d := res[0].Result.(result.Detected) //nolint:forcetypeassert
				if d.Empty() {
					t.Error("Expected detection of zero-sized types")
				}

				zerotypes := d.Sorted()
				if !slices.Contains(zerotypes, tt.want) {
					t.Errorf("Expected %q to contain %q", zerotypes, tt.want)
				}
			}
		})
	}
}
