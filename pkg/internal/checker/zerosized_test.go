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

package checker_test

import (
	"go/token"
	"go/types"
	"maps"
	"slices"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/internal/checker"
	"fillmore-labs.com/zerolint/pkg/internal/filter"
	"fillmore-labs.com/zerolint/pkg/internal/set"
)

//nolint:lll,funlen
func TestChecker_ZeroSizedType(t *testing.T) {
	t.Parallel()

	src := `
package testpkg

type EmptyStruct struct{}
type ArrayOfZero [0]int
type StructWithEmptyField struct { _ EmptyStruct }
type StructWithZeroArrayField struct { _ [0]string }
type StructTwoEmptyFields struct {
	_ EmptyStruct
	_ ArrayOfZero
}
type AliasToEmptyStruct = EmptyStruct
type EmptyStructAlias = struct{}

type NonEmptyStruct struct{ i int }
type ArrayOfOne [1]int
type StructWithNonEmptyField struct { _ NonEmptyStruct }
type AliasToNonEmptyStruct = NonEmptyStruct
type StructWithMixedFields struct {
    _ EmptyStruct
    I int
}
type ZSWithValueReceiver struct{}
func (z ZSWithValueReceiver) ValRec() {}
func (z *ZSWithValueReceiver) PtrRec() {}

type ZSWithPtrReceiverOnly struct{}
func (z *ZSWithPtrReceiverOnly) PtrRec() {}

type ZSWithEmbeddedValueReceiver struct { ZSWithValueReceiver }
type ZSWithEmbeddedPtrReceiverOnly struct { ZSWithPtrReceiverOnly }
type ZSWithNonEmbeddedValueReceiver struct { F ZSWithValueReceiver }
type ZSWithNonEmbeddedNonZero struct { F NonEmptyStruct }

type ExcludableEmptyStruct struct{}
`
	pkg := parseSource(t, "test.go", src)

	tests := [...]struct {
		name             string
		getTypeFn        func() types.Type
		setupChecker     func(c *Checker)
		wantZeroSized    bool
		wantValueMethod  bool
		wantDetectedName string // if non-empty, check if this name is in c.Detected
	}{
		{name: "basic int", getTypeFn: func() types.Type { return types.Typ[types.Int] }, wantZeroSized: false},
		{name: "pointer type (*int)", getTypeFn: func() types.Type { return types.NewPointer(types.Typ[types.Int]) }, wantZeroSized: false},
		{name: "interface type", getTypeFn: func() types.Type { return types.NewInterfaceType(nil, nil).Complete() }, wantZeroSized: false},
		{name: "slice type", getTypeFn: func() types.Type { return types.NewSlice(types.Typ[types.Int]) }, wantZeroSized: false},
		{name: "map type", getTypeFn: func() types.Type { return types.NewMap(types.Typ[types.String], types.Typ[types.Int]) }, wantZeroSized: false},

		{name: "EmptyStruct", getTypeFn: func() types.Type { return getType(t, pkg, "EmptyStruct") }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.EmptyStruct"},
		{name: "ArrayOfZero", getTypeFn: func() types.Type { return getType(t, pkg, "ArrayOfZero") }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.ArrayOfZero"},
		{name: "StructWithEmptyField", getTypeFn: func() types.Type { return getType(t, pkg, "StructWithEmptyField") }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.StructWithEmptyField"},
		{name: "StructWithZeroArrayField", getTypeFn: func() types.Type { return getType(t, pkg, "StructWithZeroArrayField") }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.StructWithZeroArrayField"},
		{name: "StructTwoEmptyFields", getTypeFn: func() types.Type { return getType(t, pkg, "StructTwoEmptyFields") }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.StructTwoEmptyFields"},
		{name: "AliasToEmptyStruct", getTypeFn: func() types.Type { return getType(t, pkg, "AliasToEmptyStruct") }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.AliasToEmptyStruct"},

		{name: "EmptyStructAlias", getTypeFn: func() types.Type { return getType(t, pkg, "EmptyStructAlias") }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.EmptyStructAlias"},
		{name: "struct{}", getTypeFn: func() types.Type { return types.Unalias(getType(t, pkg, "EmptyStructAlias")) }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "struct{}"},

		{name: "NonEmptyStruct", getTypeFn: func() types.Type { return getType(t, pkg, "NonEmptyStruct") }, wantZeroSized: false},
		{name: "ArrayOfOne", getTypeFn: func() types.Type { return getType(t, pkg, "ArrayOfOne") }, wantZeroSized: false},
		{name: "StructWithNonEmptyField", getTypeFn: func() types.Type { return getType(t, pkg, "StructWithNonEmptyField") }, wantZeroSized: false},
		{name: "AliasToNonEmptyStruct", getTypeFn: func() types.Type { return getType(t, pkg, "AliasToNonEmptyStruct") }, wantZeroSized: false},
		{name: "StructWithMixedFields", getTypeFn: func() types.Type { return getType(t, pkg, "StructWithMixedFields") }, wantZeroSized: false},

		{name: "ZSWithValueReceiver", getTypeFn: func() types.Type { return getType(t, pkg, "ZSWithValueReceiver") }, wantZeroSized: true, wantValueMethod: true, wantDetectedName: "testpkg.ZSWithValueReceiver"},
		{name: "ZSWithPtrReceiverOnly", getTypeFn: func() types.Type { return getType(t, pkg, "ZSWithPtrReceiverOnly") }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.ZSWithPtrReceiverOnly"},
		{name: "ZSWithEmbeddedValueReceiver", getTypeFn: func() types.Type { return getType(t, pkg, "ZSWithEmbeddedValueReceiver") }, wantZeroSized: true, wantValueMethod: true, wantDetectedName: "testpkg.ZSWithEmbeddedValueReceiver"},
		{name: "ZSWithEmbeddedPtrReceiverOnly", getTypeFn: func() types.Type { return getType(t, pkg, "ZSWithEmbeddedPtrReceiverOnly") }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.ZSWithEmbeddedPtrReceiverOnly"},
		{name: "ZSWithNonEmbeddedValueReceiver", getTypeFn: func() types.Type { return getType(t, pkg, "ZSWithNonEmbeddedValueReceiver") }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.ZSWithNonEmbeddedValueReceiver"},
		{name: "ZSWithNonEmbeddedNonZero", getTypeFn: func() types.Type { return getType(t, pkg, "ZSWithNonEmbeddedNonZero") }, wantZeroSized: false},

		{name: "nil", getTypeFn: func() types.Type { return nil }, wantZeroSized: false},

		{
			name:          "ExcludableEmptyStruct - excluded",
			getTypeFn:     func() types.Type { return getType(t, pkg, "ExcludableEmptyStruct") },
			setupChecker:  func(c *Checker) { c.Excludes.Add("testpkg.ExcludableEmptyStruct") },
			wantZeroSized: false, wantDetectedName: "testpkg.ExcludableEmptyStruct",
		},
		{
			name:          "ExcludableEmptyStruct - not excluded",
			getTypeFn:     func() types.Type { return getType(t, pkg, "ExcludableEmptyStruct") },
			wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.ExcludableEmptyStruct",
		},
		{
			name:      "ExcludableEmptyStruct - excluded via directive",
			getTypeFn: func() types.Type { return getType(t, pkg, "ExcludableEmptyStruct") },
			setupChecker: func(c *Checker) {
				ex, _ := getType(t, pkg, "ExcludableEmptyStruct").(*types.Named)
				c.ExcludedTypeDefs = filter.New(set.New(ex.Obj().Pos()))
			},
			wantZeroSized: false,
		},

		{
			name: "Recursive",
			getTypeFn: func() types.Type {
				// This is not a valid Go type
				name := types.NewTypeName(token.NoPos, nil, "Recursive", nil)
				recursive := types.NewNamed(name, nil, nil)
				recursive.SetUnderlying(types.NewArray(recursive, 1))

				return recursive
			},
			wantZeroSized: false,
		},
		{
			name: "Recursive2",
			getTypeFn: func() types.Type {
				// This is not a valid Go type
				name := types.NewTypeName(token.NoPos, nil, "Recursive2", nil)
				recursive := types.NewNamed(name, nil, nil)
				recursivefield := types.NewVar(token.NoPos, nil, "_", recursive)

				field := types.NewVar(token.NoPos, nil, "_", types.NewStruct(nil, nil))
				fields := slices.Repeat([]*types.Var{field}, 1_000)
				fields = append(fields, recursivefield)

				recursive.SetUnderlying(types.NewStruct(fields, nil))

				return recursive
			},
			wantZeroSized: false,
		},
		{
			name: "Big",
			getTypeFn: func() types.Type {
				field := types.NewVar(token.NoPos, nil, "_", types.NewArray(types.NewStruct(nil, nil), 1))
				for range 10 {
					field = types.NewVar(token.NoPos, nil, "_", types.NewStruct([]*types.Var{field, field}, nil))
				}

				nonEmptyField := types.NewVar(token.NoPos, nil, "_", types.NewArray(types.Typ[types.Bool], 1))
				for range 10 {
					nonEmptyField = types.NewVar(token.NoPos, nil, "_", types.NewStruct([]*types.Var{nonEmptyField}, nil))
				}

				fields := slices.Repeat([]*types.Var{field}, 1_000)
				fields = append(fields, nonEmptyField)

				return types.NewStruct(fields, nil)
			},
			wantZeroSized: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newTestChecker(t) // New checker for each test for isolation

			if tt.setupChecker != nil {
				tt.setupChecker(c)
			}

			typ := tt.getTypeFn()
			gotValueMethod, gotZeroSized := c.ZeroSizedType(typ)

			if gotZeroSized != tt.wantZeroSized {
				t.Errorf("ZeroSizedType() gotZeroSized = %v, want %v for type %s", gotZeroSized, tt.wantZeroSized, typ.String())
			}

			if gotZeroSized && gotValueMethod != tt.wantValueMethod {
				t.Errorf("ZeroSizedType() gotValueMethod = %v, want %v for type %s", gotValueMethod, tt.wantValueMethod, typ.String())
			}

			if tt.wantDetectedName != "" {
				if _, ok := c.Detected[tt.wantDetectedName]; !ok {
					detectedItems := slices.Sorted(maps.Keys(c.Detected))
					t.Errorf("ZeroSizedType() expected %q to be in Detected set, but it was not. Detected: %v", tt.wantDetectedName, detectedItems)
				}
			}

			if _, cachedZeroSized := c.ZeroSizedType(typ); cachedZeroSized != gotZeroSized {
				t.Errorf("ZeroSizedType() cachedZeroSized = %v, want %v for type %s", cachedZeroSized, gotZeroSized, typ.String())
			}
		})
	}
}

//nolint:lll,funlen
func TestChecker_ZeroSizedTypePointer(t *testing.T) {
	t.Parallel()

	src := `
package testpkg

type EmptyStruct struct{}
type NonEmptyStruct struct{ i int }

type ZSWithValueReceiver struct{}
func (z ZSWithValueReceiver) ValRec() {}

type ZSWithPtrReceiverOnly struct{}
func (z *ZSWithPtrReceiverOnly) PtrRec() {}

type ExcludableEmptyStruct struct{}

type PtrToEmptyStruct *EmptyStruct
type PtrToNonEmptyStruct *NonEmptyStruct
type PtrToZSWithValueReceiver *ZSWithValueReceiver
type PtrToZSWithPtrReceiverOnly *ZSWithPtrReceiverOnly
type PtrToExcludableEmptyStruct *ExcludableEmptyStruct
`
	pkg := parseSource(t, "test.go", src)

	tests := [...]struct {
		name             string
		getTypeFn        func() types.Type
		setupChecker     func(c *Checker)
		wantElemTypeName string
		wantZeroSized    bool
		wantValueMethod  bool
	}{
		{
			name:          "nil type",
			getTypeFn:     func() types.Type { return nil },
			wantZeroSized: false,
		},
		{
			name:          "non-pointer type (EmptyStruct)",
			getTypeFn:     func() types.Type { return getType(t, pkg, "EmptyStruct") },
			wantZeroSized: false,
		},
		{
			name:          "pointer to non-zero-sized type (*NonEmptyStruct)",
			getTypeFn:     func() types.Type { return getType(t, pkg, "PtrToNonEmptyStruct") },
			wantZeroSized: false,
		},
		{
			name:          "pointer to non-zero-sized basic type (*int)",
			getTypeFn:     func() types.Type { return types.NewPointer(types.Typ[types.Int]) },
			wantZeroSized: false,
		},
		{
			name:             "pointer to zero-sized type (*EmptyStruct)",
			getTypeFn:        func() types.Type { return getType(t, pkg, "PtrToEmptyStruct") },
			wantElemTypeName: "testpkg.EmptyStruct", wantZeroSized: true, wantValueMethod: false,
		},
		{
			name:             "pointer to ZSWithValueReceiver",
			getTypeFn:        func() types.Type { return getType(t, pkg, "PtrToZSWithValueReceiver") },
			wantElemTypeName: "testpkg.ZSWithValueReceiver", wantZeroSized: true, wantValueMethod: true,
		},
		{
			name:             "pointer to ZSWithPtrReceiverOnly",
			getTypeFn:        func() types.Type { return getType(t, pkg, "PtrToZSWithPtrReceiverOnly") },
			wantElemTypeName: "testpkg.ZSWithPtrReceiverOnly", wantZeroSized: true, wantValueMethod: false,
		},
		{
			name:          "pointer to ExcludableEmptyStruct - excluded",
			getTypeFn:     func() types.Type { return getType(t, pkg, "PtrToExcludableEmptyStruct") },
			setupChecker:  func(c *Checker) { c.Excludes.Add("testpkg.ExcludableEmptyStruct") },
			wantZeroSized: false,
		},
		{
			name:             "pointer to ExcludableEmptyStruct - not excluded",
			getTypeFn:        func() types.Type { return getType(t, pkg, "PtrToExcludableEmptyStruct") },
			wantElemTypeName: "testpkg.ExcludableEmptyStruct", wantZeroSized: true, wantValueMethod: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newTestChecker(t) // New checker for each test for isolation

			if tt.setupChecker != nil {
				tt.setupChecker(c)
			}

			typToTest := tt.getTypeFn()
			elem, valueMethod, zeroSized := c.ZeroSizedTypePointer(typToTest)

			if zeroSized != tt.wantZeroSized {
				t.Errorf("ZeroSizedTypePointer() zeroSized = %v, want %v for input type %s", zeroSized, tt.wantZeroSized, typToTest)
			}

			if zeroSized {
				if elem == nil {
					t.Fatalf("ZeroSizedTypePointer() elem is nil, want non-nil (%s)", tt.wantElemTypeName)
				}

				if elem.String() != tt.wantElemTypeName {
					t.Errorf("ZeroSizedTypePointer() elem.String() = %q, want %q", elem.String(), tt.wantElemTypeName)
				}

				if valueMethod != tt.wantValueMethod {
					t.Errorf("ZeroSizedTypePointer() valueMethod = %v, want %v for input type %s", valueMethod, tt.wantValueMethod, typToTest)
				}
			}
		})
	}
}
