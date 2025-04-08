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
	"go/types"
	"testing"
)

func TestChecker_TypesInfo(t *testing.T) {
	t.Parallel()

	info, pkg, fset, astFile := parseSource(t, "main.go", "package main")
	c := newTestChecker(t, info, pkg, fset, astFile)

	if c.TypesInfo() != info {
		t.Errorf("TypesInfo() = %v, want %v", c.TypesInfo(), info)
	}
}

func TestChecker_Print(t *testing.T) {
	t.Parallel()

	info, pkg, fset, astFile := parseSource(t, "main.go", "package main")
	c := newTestChecker(t, info, pkg, fset, astFile)

	err := c.Print(astFile)
	if err != nil {
		t.Errorf("Got error %v printing test file", err)
	}
}

func TestChecker_IgnoreType(t *testing.T) {
	t.Parallel()

	info, pkg, fset, astFile := parseSource(t, "main.go", "package main\ntype MyType struct{}")
	c := newTestChecker(t, info, pkg, fset, astFile)

	typeName, _ := getType(t, pkg, "MyType").(*types.Named)

	// Before ignoring, MyType (struct{}) should be considered zero-sized.
	if _, zS := c.ZeroSizedType(typeName); !zS {
		t.Error("MyType (struct{}) was not considered zero-sized before IgnoreType call, but it should be.")
	}

	c.IgnoreType(typeName.Obj())

	// After ignoring, MyType should no longer be considered zero-sized by ZeroSizedType.
	if _, zS := c.ZeroSizedType(typeName); zS {
		t.Error("IgnoreType() did not cause MyType to be ignored by ZeroSizedType.")
	}
}
