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

// operandInfo holds type information for comparison operands.
type operandInfo struct {
	zeroSizedPointer bool
	elem             types.Type
}

type comparison struct {
	left, right operandInfo
}

// visitCmp analyzes comparison expressions (x == y, x != y, errors.Is(x, y)) for comparisons
// involving pointers to zero-sized types.
func (v *Visitor) visitCmp(n ast.Node, x, y ast.Expr) bool {
	var p comparison
	var ok bool
	if p.left, ok = v.operandInfo(x); !ok {
		return true
	}
	if p.right, ok = v.operandInfo(y); !ok {
		return true
	}

	var message string
	switch {
	case p.left.zeroSizedPointer && p.right.zeroSizedPointer:
		message = p.comparisonMessage()

	case p.left.zeroSizedPointer:
		message = p.comparisonMessageLeft()

	case p.right.zeroSizedPointer:
		message = p.comparisonMessageRight()

	default:
		return true
	}

	v.report(n, message, nil)

	// no fixes, so dive deeper.
	return true
}

// operandInfo extracts type information about comparison operands,
// identifying whether they are pointers to zero-sized types or interfaces.
func (v *Visitor) operandInfo(x ast.Expr) (operandInfo, bool) {
	tv := v.Pass.TypesInfo.Types[x]
	if tv.IsNil() {
		return operandInfo{}, false // comparisons to nil are not flagged
	}

	switch t := tv.Type.Underlying().(type) {
	case *types.Pointer:
		if !v.zeroSizedType(t.Elem()) {
			return operandInfo{}, false
		}

		return operandInfo{zeroSizedPointer: true, elem: t.Elem()}, true // comparisons with a pointer to zero-size vriable

	case *types.Interface:
		return operandInfo{elem: tv.Type}, true // comparisons with an interface

	default:
		return operandInfo{}, false // other comparisons
	}
}

// comparisonMessage generates an appropriate diagnostic message for comparing
// two pointers to zero-sized types.
func (c comparison) comparisonMessage() string {
	if types.Identical(c.left.elem, c.right.elem) { // types.Identical ignores aliases
		return fmt.Sprintf("comparison of pointers to zero-size variables of type %q (ZL01)", c.left.elem)
	}

	return fmt.Sprintf("comparison of pointers to zero-size variables of types %q and %q (ZL01)",
		c.left.elem, c.right.elem)
}

// comparisonMessageLeft generates a diagnostic message for pointer-to-interface comparison.
func (c comparison) comparisonMessageLeft() string {
	return fmt.Sprintf("comparison of pointer to zero-size variable of type %q with interface of type %q (ZL02)",
		c.left.elem, c.right.elem)
}

// comparisonMessageLeft generates a diagnostic message for interface-to-pointer comparison.
func (c comparison) comparisonMessageRight() string {
	return fmt.Sprintf("comparison of pointer to zero-size variable of type %q with interface of type %q (ZL02)",
		c.right.elem, c.left.elem)
}
