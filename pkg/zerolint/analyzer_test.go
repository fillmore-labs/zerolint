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

package zerolint_test

import (
	"bytes"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"fillmore-labs.com/zerolint/pkg/internal/excludes"
	. "fillmore-labs.com/zerolint/pkg/zerolint"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

func TestAnalyzer(t *testing.T) { //nolint:funlen
	t.Parallel()

	dir := analysistest.TestData()

	excludedTypeNames, err := excludes.ReadExcludes(os.DirFS(dir), "excluded.txt")
	if err != nil {
		t.Fatalf("Can't find excludes file: %v", err)
	}

	tests := []struct {
		name    string
		options Options
		flags   map[string]string
		want    []string
		pkg     string
	}{
		{
			name: "basic with zerotrace",
			options: Options{
				WithLevel(level.Basic),
				WithExcludes(excludedTypeNames),
				WithZeroTrace(true),
				WithGenerated(true),
				WithExcludeComments(true),
				WithRegex(regexp.MustCompile(`^test/`)),
			},
			want: []string{"test/basic.aliasError", "test/basic.myError"},
			pkg:  "test/basic",
		},
		{
			name: "basic with zerotrace via flags",
			options: Options{
				WithFlags(true),
			},
			flags: map[string]string{
				"level":     "basic",
				"excluded":  dir + "/excluded.txt",
				"zerotrace": "true",
				"generated": "true",
				"match":     `^test/`,
			},
			want: []string{"test/basic.aliasError", "test/basic.myError"},
			pkg:  "test/basic",
		},
		{
			name: "no exclusions",
			options: Options{
				WithLevel(level.Full),
				WithZeroTrace(true),
				WithExcludeComments(false),
			},
			want: []string{"test/noexclude.excludedError"},
			pkg:  "test/noexclude",
		},
	}
	for _, tt := range tests {
		var buf bytes.Buffer
		logOpt := WithLogger(log.New(&buf, "", 0)) // Capture logs
		a := New(logOpt, tt.options)

		for key, value := range tt.flags {
			if err := a.Flags.Set(key, value); err != nil {
				t.Fatalf("Can't set flag %s=%s: %v", key, value, err)
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			analysistest.Run(t, dir, a, tt.pkg)

			// Assert log output
			got := buf.String()
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("expected log to contain %q, got:\n%s", want, got)
				}
			}
		})
	}
}
