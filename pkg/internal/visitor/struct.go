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

package visitor

import (
	"fmt"
	"go/ast"
	"go/types"
)

// visitStructType analyzes struct type declarations for fields or embedded types
// that are pointers to zero-sized types.
// If `level` is false, it only checks embedded types.
func (v *Visitor) visitStructType(n *ast.StructType) bool {
	v.checkFieldList(n.Fields, v.level < 2, msgStruct{}) //nolint:mnd

	return true
}

// msgStruct implements msgFormatter for struct fields, providing appropriate messages
// for diagnostics related to struct fields pointing to zero-sized types.
type msgStruct struct{}

// ZeroMsg returns a message for embedded pointer fields to zero-sized types.
func (msgStruct) ZeroMsg(t types.Type) (category, string) {
	return catStructEmbedded, fmt.Sprintf("embedded pointer to zero-sized type %q", t)
}

// SingularMsg returns a message for a single named struct field pointing to a zero-sized type.
func (msgStruct) SingularMsg(name string, t types.Type) (category, string) {
	return catStructField, fmt.Sprintf("field %q points to zero-sized type %q", name, t)
}

// PluralMsg returns a message for multiple struct fields pointing to a zero-sized type.
func (msgStruct) PluralMsg(names string, t types.Type) (category, string) {
	return catStructField, fmt.Sprintf("fields %s point to zero-sized type %q", names, t)
}
