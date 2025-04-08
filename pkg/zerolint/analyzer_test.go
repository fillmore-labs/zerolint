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

package zerolint_test

import (
	"bytes"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"fillmore-labs.com/zerolint/pkg/internal/excludes"
	. "fillmore-labs.com/zerolint/pkg/zerolint"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) { //nolint:funlen
	t.Parallel()

	dir := analysistest.TestData()

	ex, err := excludes.ReadExcludes(os.DirFS(dir), "excluded.txt")
	if err != nil {
		t.Fatalf("Can't find excludes file: %v", err)
	}

	tests := []struct {
		name    string
		options []Option
		flags   map[string]string
		want    []string
		pkg     string
	}{
		{
			name: "basic with zerotrace",
			options: []Option{
				WithLevel(level.Default),
				WithExcludes(ex),
				WithZeroTrace(true),
				WithGenerated(true),
				WithRegex(regexp.MustCompile(`^go\.test/`)),
			},
			want: []string{"go.test/basic.aliasError", "go.test/basic.myError"},
			pkg:  "go.test/basic",
		},
		{
			name: "basic with zerotrace via flags",
			options: []Option{
				WithFlags(true),
			},
			flags: map[string]string{
				"level":     "default",
				"excluded":  dir + "/excluded.txt",
				"zerotrace": "true",
				"generated": "true",
				"match":     `^go\.test/`,
			},
			want: []string{"go.test/basic.aliasError", "go.test/basic.myError"},
			pkg:  "go.test/basic",
		},
	}
	for _, tt := range tests {
		var buf bytes.Buffer // Capture logs
		alloptions := []Option{WithLogger(log.New(&buf, "", 0))}
		alloptions = append(alloptions, tt.options...)

		a := New(alloptions...)

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
