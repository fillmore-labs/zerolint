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
	"errors"
	"io/fs"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	. "fillmore-labs.com/zerolint/pkg/zerolint"
)

type ignoreTestErrors struct{}

func (ignoreTestErrors) Errorf(_ string, _ ...any) {}

func TestAnalyzerWithExcludedFlag(t *testing.T) {
	t.Parallel()

	a := New(WithFlags(true))
	a.RunDespiteErrors = true

	if err := a.Flags.Set("excluded", "/nonexistent"); err != nil {
		t.Skipf("Can't set excluded: %v", err)
	}

	dir := analysistest.TestData()
	result := analysistest.Run(ignoreTestErrors{}, dir, a, "test/none")

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	if err := result[0].Action.Err; !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("wanted %v, got: %v", fs.ErrNotExist, err)
	}
}
