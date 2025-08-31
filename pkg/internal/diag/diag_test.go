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

package diag_test

import (
	"testing"
)

func TestDiag_TypesInfo(t *testing.T) {
	t.Parallel()

	info, pkg, fset, astFile := parseSource(t, "main.go", "package main")

	if d := newTestDiag(t, info, pkg, fset, astFile); d.TypesInfo() != info {
		t.Errorf("TypesInfo() = %v, want %v", d.TypesInfo(), info)
	}
}
