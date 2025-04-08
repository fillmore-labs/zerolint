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
	"go/ast"
	"go/types"

	"fillmore-labs.com/zerolint/pkg/internal/checker"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

// visitStructType analyzes struct type declarations for fields or embedded types
// that are pointers to zero-sized types.
// If the lint level is `Default` (i.e., `v.level.Below(level.Extended)` is true), it only checks embedded types.
func (v *Visitor) visitStructType(n *ast.StructType) bool {
	v.checkFieldList(n.Fields, v.level.Below(level.Extended), msgStruct{})

	return true
}

// msgStruct implements msgFormatter for struct fields, providing appropriate messages
// for diagnostics related to struct fields pointing to zero-sized types.
type msgStruct struct{}

// zeroMsg returns a message for embedded pointer fields to zero-sized types.
func (msgStruct) zeroMsg(typ types.Type, valueMethod bool) checker.CategorizedMessage {
	return msgFormatf(catStructEmbedded, valueMethod, "embedded pointer to zero-sized type %q", typ)
}

// singularMsg returns a message for a single named struct field pointing to a zero-sized type.
func (msgStruct) singularMsg(typ types.Type, valueMethod bool, name string) checker.CategorizedMessage {
	return msgFormatf(catStructField, valueMethod, "field %q points to zero-sized type %q", name, typ)
}

// pluralMsg returns a message for multiple struct fields pointing to a zero-sized type.
func (msgStruct) pluralMsg(typ types.Type, valueMethod bool, names string) checker.CategorizedMessage {
	return msgFormatf(catStructField, valueMethod, "fields %s point to zero-sized type %q", names, typ)
}
