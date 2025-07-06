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

package analyzer

import (
	"go/ast"

	"fillmore-labs.com/zerolint/pkg/internal/analyzer/msg"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

// visitStructType analyzes struct type declarations for fields or embedded types
// that are pointers to zero-sized types.
// If the lint level is `Default` (i.e., `v.level.Below(level.Extended)` is true), it only checks embedded types.
func (v *Visitor) visitStructType(n *ast.StructType) bool {
	v.checkFieldList(n.Fields, v.Level.Below(level.Extended), msg.Struct{})

	return true
}
