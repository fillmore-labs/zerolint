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

package base_test

import (
	"go/types"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/internal/base"
	"fillmore-labs.com/zerolint/pkg/internal/set"
)

//nolint:lll,funlen,paralleltest
func TestVisitor_ZeroSizedType(t *testing.T) {
	src := `
package testpkg

import "runtime"

type EmptyStruct struct{}
type ArrayOfZero [0]int
type StructWithEmptyField struct { _ EmptyStruct }
type StructWithZeroArrayField struct { _ [0]string }
type StructTwoEmptyFields struct {
	_ EmptyStruct
	_ ArrayOfZero
}
type AliasToEmptyStruct = EmptyStruct

type NonEmptyStruct struct{ i int }
type ArrayOfOne [1]int
type StructWithNonEmptyField struct { _ NonEmptyStruct }
type AliasToNonEmptyStruct = NonEmptyStruct
type StructWithMixedFields struct {
    _ EmptyStruct
    I int
}
type RuntimeFuncAlias = runtime.Func
type _Ctype_struct struct{}

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
type IgnoredType struct{} // Will be struct{}
`
	info, pkg, fset, astFile := parseSource(t, src)

	var ignoredType *types.TypeName

	obj := pkg.Scope().Lookup("IgnoredType")
	if tn, ok := obj.(*types.TypeName); ok {
		ignoredType = tn
	} else {
		t.Fatal("IgnoredType not found or not a TypeName for Pos setup")
	}

	tests := []struct {
		name             string
		getTypeFn        func() types.Type
		setupVisitor     func(v *Base)
		wantZeroSized    bool
		wantValueMethod  bool
		wantDetectedName string // if non-empty, check if this name is in v.Detected
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

		{name: "NonEmptyStruct", getTypeFn: func() types.Type { return getType(t, pkg, "NonEmptyStruct") }, wantZeroSized: false},
		{name: "ArrayOfOne", getTypeFn: func() types.Type { return getType(t, pkg, "ArrayOfOne") }, wantZeroSized: false},
		{name: "StructWithNonEmptyField", getTypeFn: func() types.Type { return getType(t, pkg, "StructWithNonEmptyField") }, wantZeroSized: false},
		{name: "AliasToNonEmptyStruct", getTypeFn: func() types.Type { return getType(t, pkg, "AliasToNonEmptyStruct") }, wantZeroSized: false},
		{name: "StructWithMixedFields", getTypeFn: func() types.Type { return getType(t, pkg, "StructWithMixedFields") }, wantZeroSized: false},
		{name: "RuntimeFunc", getTypeFn: func() types.Type { return types.Unalias(getType(t, pkg, "RuntimeFuncAlias")) }, wantZeroSized: false},
		{name: "Ctype", getTypeFn: func() types.Type { return getType(t, pkg, "_Ctype_struct") }, wantZeroSized: false},

		{name: "ZSWithValueReceiver", getTypeFn: func() types.Type { return getType(t, pkg, "ZSWithValueReceiver") }, wantZeroSized: true, wantValueMethod: true, wantDetectedName: "testpkg.ZSWithValueReceiver"},
		{name: "ZSWithPtrReceiverOnly", getTypeFn: func() types.Type { return getType(t, pkg, "ZSWithPtrReceiverOnly") }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.ZSWithPtrReceiverOnly"},
		{name: "ZSWithEmbeddedValueReceiver", getTypeFn: func() types.Type { return getType(t, pkg, "ZSWithEmbeddedValueReceiver") }, wantZeroSized: true, wantValueMethod: true, wantDetectedName: "testpkg.ZSWithEmbeddedValueReceiver"},
		{name: "ZSWithEmbeddedPtrReceiverOnly", getTypeFn: func() types.Type { return getType(t, pkg, "ZSWithEmbeddedPtrReceiverOnly") }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.ZSWithEmbeddedPtrReceiverOnly"},
		{name: "ZSWithNonEmbeddedValueReceiver", getTypeFn: func() types.Type { return getType(t, pkg, "ZSWithNonEmbeddedValueReceiver") }, wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.ZSWithNonEmbeddedValueReceiver"},
		{name: "ZSWithNonEmbeddedNonZero", getTypeFn: func() types.Type { return getType(t, pkg, "ZSWithNonEmbeddedNonZero") }, wantZeroSized: false},

		{
			name:          "ExcludableEmptyStruct - excluded",
			getTypeFn:     func() types.Type { return getType(t, pkg, "ExcludableEmptyStruct") },
			setupVisitor:  func(v *Base) { v.Excludes.Insert("testpkg.ExcludableEmptyStruct") },
			wantZeroSized: false, wantDetectedName: "testpkg.ExcludableEmptyStruct",
		},
		{
			name:          "ExcludableEmptyStruct - not excluded",
			getTypeFn:     func() types.Type { return getType(t, pkg, "ExcludableEmptyStruct") },
			wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.ExcludableEmptyStruct",
		},
		{
			name:          "IgnoredType - marked ignored",
			getTypeFn:     func() types.Type { return getType(t, pkg, "IgnoredType") },
			setupVisitor:  func(v *Base) { v.IgnoreType(ignoredType) },
			wantZeroSized: false,
		},
		{
			name:          "IgnoredType - not marked ignored",
			getTypeFn:     func() types.Type { return getType(t, pkg, "IgnoredType") },
			wantZeroSized: true, wantValueMethod: false, wantDetectedName: "testpkg.IgnoredType",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := newTestVisitor(t, info, pkg, fset, astFile) // New visitor for each test for isolation
			if tt.setupVisitor != nil {
				tt.setupVisitor(v)
			}

			typ := tt.getTypeFn()
			gotValueMethod, gotZeroSized := v.ZeroSizedType(typ)

			if gotZeroSized != tt.wantZeroSized {
				t.Errorf("ZeroSizedType() gotZeroSized = %v, want %v for type %s", gotZeroSized, tt.wantZeroSized, typ.String())
			}

			if gotZeroSized && gotValueMethod != tt.wantValueMethod {
				t.Errorf("ZeroSizedType() gotValueMethod = %v, want %v for type %s", gotValueMethod, tt.wantValueMethod, typ.String())
			}

			if tt.wantDetectedName != "" && !v.Detected.Has(tt.wantDetectedName) {
				detectedItems := set.Sorted(v.Detected)
				t.Errorf("ZeroSizedType() expected %q to be in Detected set, but it was not. Detected: %v", tt.wantDetectedName, detectedItems)
			}
		})
	}
}

//nolint:lll,funlen,paralleltest
func TestVisitor_ZeroSizedTypePointer(t *testing.T) {
	src := `
package testpkg

type EmptyStruct struct{}
type NonEmptyStruct struct{ i int }

type ZSWithValueReceiver struct{}
func (z ZSWithValueReceiver) ValRec() {}

type ZSWithPtrReceiverOnly struct{}
func (z *ZSWithPtrReceiverOnly) PtrRec() {}

type ExcludableEmptyStruct struct{}
type IgnoredType struct{} // Will be struct{}

type PtrToEmptyStruct *EmptyStruct
type PtrToNonEmptyStruct *NonEmptyStruct
type PtrToZSWithValueReceiver *ZSWithValueReceiver
type PtrToZSWithPtrReceiverOnly *ZSWithPtrReceiverOnly
type PtrToExcludableEmptyStruct *ExcludableEmptyStruct
type PtrToIgnoredType *IgnoredType
`
	info, pkg, fset, astFile := parseSource(t, src)

	var ignoredType *types.TypeName

	objIgnored := pkg.Scope().Lookup("IgnoredType")
	if tn, ok := objIgnored.(*types.TypeName); ok {
		ignoredType = tn
	} else {
		t.Fatal("IgnoredType not found or not a TypeName for Pos setup")
	}

	tests := []struct {
		name             string
		getTypeFn        func() types.Type
		setupVisitor     func(v *Base)
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
			setupVisitor:  func(v *Base) { v.Excludes.Insert("testpkg.ExcludableEmptyStruct") },
			wantZeroSized: false,
		},
		{
			name:             "pointer to ExcludableEmptyStruct - not excluded",
			getTypeFn:        func() types.Type { return getType(t, pkg, "PtrToExcludableEmptyStruct") },
			wantElemTypeName: "testpkg.ExcludableEmptyStruct", wantZeroSized: true, wantValueMethod: false,
		},
		{
			name:          "pointer to IgnoredType - marked ignored",
			getTypeFn:     func() types.Type { return getType(t, pkg, "PtrToIgnoredType") },
			setupVisitor:  func(v *Base) { v.IgnoreType(ignoredType) },
			wantZeroSized: false,
		},
		{
			name:             "pointer to IgnoredType - not marked ignored",
			getTypeFn:        func() types.Type { return getType(t, pkg, "PtrToIgnoredType") },
			wantElemTypeName: "testpkg.IgnoredType", wantZeroSized: true, wantValueMethod: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := newTestVisitor(t, info, pkg, fset, astFile) // New visitor for each test for isolation
			if tt.setupVisitor != nil {
				tt.setupVisitor(v)
			}

			typToTest := tt.getTypeFn()
			elem, valueMethod, zeroSized := v.ZeroSizedTypePointer(typToTest)

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
