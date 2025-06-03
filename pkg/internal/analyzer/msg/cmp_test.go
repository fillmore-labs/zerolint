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

package msg_test

import (
	"go/token"
	"go/types"
	"strings"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/internal/analyzer/msg"
)

func TestComparisonMessage(t *testing.T) {
	t.Parallel()

	structType := types.NewStruct(nil, nil)
	arrayType := types.NewArray(types.Typ[types.Int], 0)
	pkg := types.NewPackage("test", "test")
	namedStructType := types.NewNamed(types.NewTypeName(token.NoPos, pkg, "MyStruct", nil), structType, nil)
	anotherNamedStructType := types.NewNamed(types.NewTypeName(token.NoPos, pkg, "MyOtherStruct", nil), structType, nil)

	testCases := []struct {
		name        string
		left        types.Type
		right       types.Type
		valueMethod bool
		want        string
	}{
		{
			name:        "identical anonymous structs",
			left:        structType,
			right:       structType,
			valueMethod: false,
			want:        `comparison of pointers to zero-size type "struct{}"`,
		},
		{
			name:        "different zero-size types",
			left:        structType,
			right:       arrayType,
			valueMethod: false,
			want:        `comparison of pointers to zero-size types "struct{}" and "[0]int"`,
		},
		{
			name:        "identical named structs",
			left:        namedStructType,
			right:       namedStructType,
			valueMethod: true,
			want:        `comparison of pointers to zero-size type "test.MyStruct"`,
		},
		{
			name:        "different named structs",
			left:        namedStructType,
			right:       anotherNamedStructType,
			valueMethod: false,
			want:        `comparison of pointers to zero-size types "test.MyStruct" and "test.MyOtherStruct"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := ComparisonMessage(tc.left, tc.right, tc.valueMethod)

			if !strings.Contains(got.Message, tc.want) {
				t.Errorf("ComparisonMessage() returned %q, does not contain %q", got.Message, tc.want)
			}
		})
	}
}

func TestComparisonMessagePointerInterface(t *testing.T) {
	t.Parallel()

	structType := types.NewStruct(nil, nil)
	pkg := types.NewPackage("test", "test")
	namedStructType := types.NewNamed(types.NewTypeName(token.NoPos, pkg, "MyStruct", nil), structType, nil)
	errorInterface := types.Universe.Lookup("error").Type()
	emptyInterface := types.NewInterfaceType(nil, nil)

	testCases := []struct {
		name        string
		elemOp      types.Type
		interfaceOp types.Type
		valueMethod bool
		want        string
	}{
		{
			name:        "with error interface",
			elemOp:      structType,
			interfaceOp: errorInterface,
			valueMethod: false,
			want:        `comparison of pointer to zero-size type "struct{}" with error interface`,
		},
		{
			name:        "with empty interface",
			elemOp:      structType,
			interfaceOp: emptyInterface,
			valueMethod: true,
			want:        `comparison of pointer to zero-size type "struct{}" with interface of type "interface{}"`,
		},
		{
			name:        "named type with error interface",
			elemOp:      namedStructType,
			interfaceOp: errorInterface,
			valueMethod: false,
			want:        `comparison of pointer to zero-size type "test.MyStruct" with error interface`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := ComparisonMessagePointerInterface(tc.elemOp, tc.interfaceOp, tc.valueMethod)

			if !strings.Contains(got.Message, tc.want) {
				t.Errorf("ComparisonMessagePointerInterface() returned %q, does not contain %q", got.Message, tc.want)
			}
		})
	}
}
