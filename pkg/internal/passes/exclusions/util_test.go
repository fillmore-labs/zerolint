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

package exclusions_test

import (
	"go/token"
	"go/types"
	"strconv"
	"testing"

	"golang.org/x/tools/go/analysis"

	. "fillmore-labs.com/zerolint/pkg/internal/passes/exclusions"
)

type myFact struct{ value int }

func (*myFact) AFact() {}

type otherFact struct{ value int }

func (*otherFact) AFact() {}

func TestAllFacts(t *testing.T) {
	t.Parallel()

	pkg := types.NewPackage("test", "test")

	facts := make([]analysis.ObjectFact, 0, 4)

	for i := range 4 {
		fact := analysis.ObjectFact{
			Object: types.NewTypeName(token.NoPos, pkg, "MyType"+strconv.Itoa(i+1), nil),
		}
		switch i {
		case 1:
			fact.Fact = &otherFact{value: i + 1}

		default:
			fact.Fact = &myFact{value: i + 1}
		}

		facts = append(facts, fact)
	}

	found := make(map[string]int, 2)

	for obj, fact := range AllFacts[*myFact](facts) {
		found[obj.Name()] = fact.value
		if len(found) == 2 {
			break
		}
	}

	if len(found) != 2 {
		t.Errorf("Expected to find 2 facts, but got %d", len(found))
	}

	if val, ok := found["MyType1"]; !ok || val != 1 {
		t.Errorf("Expected MyType1 with value 1, got %d (found: %v)", val, ok)
	}

	if val, ok := found["MyType3"]; !ok || val != 3 {
		t.Errorf("Expected MyType3 with value 3, got %d (found: %v)", val, ok)
	}

	if _, ok := found["MyType2"]; ok {
		t.Error("Expected not to find MyType2, but it was present")
	}
}
