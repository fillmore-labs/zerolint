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

package analyzer_test

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"fillmore-labs.com/zerolint/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) { //nolint:tparallel
	t.Parallel()

	dir := analysistest.TestData()

	type args struct {
		excludes string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"basic", args{dir + "/excluded.txt"}, "go.test/basic.myError"},
	}
	for _, tt := range tests { //nolint:paralleltest
		var buf bytes.Buffer

		analyzer.Excludes = tt.args.excludes
		analyzer.Logger = log.New(&buf, "", 0)
		analyzer.ZeroTrace = true

		a := analyzer.Analyzer

		t.Run(tt.name, func(t *testing.T) {
			analysistest.Run(t, dir, a, "go.test/basic")
		})

		if got := buf.String(); !strings.Contains(got, tt.want) {
			t.Errorf("expected log to contain %s, got:\n%s", tt.want, got)
		}
	}
}
