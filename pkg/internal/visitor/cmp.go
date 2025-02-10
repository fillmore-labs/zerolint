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

// comparisonInfo holds type information for comparison operands.
type comparisonInfo struct {
	elem        types.Type
	zeroPointer bool
}

type comparison struct {
	left, right comparisonInfo
}

// visitCmp checks comparisons like x == y, x != y and errors.Is(x, y).
func (v *visitorInternal) visitCmp(n ast.Node, x, y ast.Expr) bool {
	var p comparison
	var isNil bool
	if p.left, isNil = v.getComparisonInfo(x); isNil {
		return true
	}
	if p.right, isNil = v.getComparisonInfo(y); isNil {
		return true
	}

	var message string
	switch {
	case p.left.zeroPointer && p.right.zeroPointer:
		message = p.comparisonMessage()

	case p.left.zeroPointer:
		message = p.comparisonMessageLeft()

	case p.right.zeroPointer:
		message = p.comparisonMessageRight()

	default:
		return true
	}

	v.report(n, message, nil)

	// no fixes, so dive deeper.
	return true
}

// getComparisonInfo extracts relevant type information for comparison.
func (v *visitorInternal) getComparisonInfo(x ast.Expr) (info comparisonInfo, isNil bool) {
	tv := v.Pass.TypesInfo.Types[x]
	if tv.IsNil() { // nil
		return info, true
	}

	t := tv.Type
	underlying, ok := t.Underlying().(*types.Pointer)
	if !ok { // not a pointer
		info.elem = t

		return info, false
	}

	elem := underlying.Elem()
	if !v.zeroSizedType(elem) { // not a pointer to a zero-sized variable
		info.elem = t

		return info, false
	}

	// pointer to a zero-sized variable
	info.zeroPointer = true
	info.elem = elem

	return info, false
}

// comparisonMessage determines the appropriate message for comparison of two pointers.
func (c comparison) comparisonMessage() string {
	if c.left.elem == c.right.elem {
		return fmt.Sprintf("comparison of pointers to zero-size variables of type %q", c.left.elem)
	}

	return fmt.Sprintf("comparison of pointers to zero-size variables of types %q and %q", c.left.elem, c.right.elem)
}

const comparisonMessage = "comparison of pointer to zero-size variable of type %q with variable of type %q"

// comparisonMessageLeft determines the appropriate message for pointer comparison.
func (c comparison) comparisonMessageLeft() string {
	return fmt.Sprintf(comparisonMessage, c.left.elem, c.right.elem)
}

// comparisonMessageRight determines the appropriate message for pointer comparison.
func (c comparison) comparisonMessageRight() string {
	return fmt.Sprintf(comparisonMessage, c.right.elem, c.left.elem)
}
