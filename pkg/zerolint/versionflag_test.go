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
	"os"
	"os/exec"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/zerolint"
	"golang.org/x/tools/go/analysis/analysistest"
)

const versionFlagTest = "VERSION_FLAG_TEST"

func TestVFlag(t *testing.T) { //nolint:paralleltest
	if os.Getenv(versionFlagTest) != "1" {
		return
	}

	a := New(WithFlags(true))

	if err := a.Flags.Set("V", "true"); err != nil {
		t.Fatalf("Can't set V: %v", err)
	}

	dir := analysistest.TestData()
	analysistest.Run(t, dir, a, "go.test/none")

	t.Fatal("Analyzer should exit")
}

func TestAnalyzerWithVFlag(t *testing.T) {
	t.Parallel()

	// See https://go.dev/talks/2014/testing.slide#23
	cmd := exec.Command(os.Args[0], "-test.run=TestVFlag") //nolint:gosec
	cmd.Env = append(os.Environ(), versionFlagTest+"=1")

	err := cmd.Run()
	if err != nil {
		t.Fatalf("process ran with err %v, want exit status 0", err)
	}
}
