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

//go:build go1.24

package zerolint_test

import (
	"bytes"
	"io"
	"log"
	"log/slog"
	"regexp"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/zerolint"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

func TestOptions(t *testing.T) {
	t.Parallel()

	opts := Options{
		WithExcludeComments(true),
		WithExcludes([]string{"exclude1", "exclude2"}),
		WithFlags(false),
		WithGenerated(false),
		WithLevel(level.Basic),
		WithLogger(log.New(io.Discard, "test:", 0)),
		WithRegex(regexp.MustCompile("^.*$")),
		WithZeroTrace(true),
		Options{},
	}

	var buf bytes.Buffer
	h := slog.NewTextHandler(&buf, nil)
	l := slog.New(h)

	l.InfoContext(t.Context(), "test", "options", opts)

	got := buf.String()
	if len(got) == 0 {
		t.Errorf("Expected non-empty log, got: %v", got)
	}
}
