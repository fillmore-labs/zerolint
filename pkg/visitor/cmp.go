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
	elem          types.Type
	isZeroPointer bool
	isInterface   bool
}

// visitCmp checks comparisons like x == y, x != y and errors.Is(x, y).
func (v Visitor) visitCmp(n ast.Node, x, y ast.Expr) bool {
	var p [2]comparisonInfo
	for i, z := range []ast.Expr{x, y} {
		t := v.TypesInfo.Types[z]
		if t.IsNil() {
			return true
		}
		p[i] = v.getComparisonInfo(t.Type)
	}

	var message string
	switch {
	case p[0].isZeroPointer && p[1].isZeroPointer:
		message = comparisonMessage(p[0].elem, p[1].elem)

	case p[0].isZeroPointer:
		message = comparisonIMessage(p[0].elem, p[1].elem, p[1].isInterface)

	case p[1].isZeroPointer:
		message = comparisonIMessage(p[1].elem, p[0].elem, p[0].isInterface)

	default:
		return true
	}

	v.report(n, message, nil)

	// no fixes, so dive deeper.
	return true
}

// getComparisonInfo extracts relevant type information for comparison.
func (v Visitor) getComparisonInfo(t types.Type) comparisonInfo {
	var info comparisonInfo

	switch underlying := t.Underlying().(type) {
	case *types.Pointer:
		info.elem = underlying.Elem()
		info.isZeroPointer = v.isZeroSizedType(info.elem)

	case *types.Interface:
		info.elem = t
		info.isInterface = true

	default:
		info.elem = t
	}

	return info
}

// comparisonMessage determines the appropriate message for comparison of two pointers.
func comparisonMessage(xType, yType types.Type) string {
	if types.Identical(xType, yType) {
		return fmt.Sprintf("comparison of pointers to zero-size variables of type %q", xType)
	}

	return fmt.Sprintf("comparison of pointers to zero-size variables of types %q and %q", xType, yType)
}

// comparisonIMessage determines the appropriate message for pointer comparison.
func comparisonIMessage(zType, withType types.Type, isInterface bool) string {
	var isIf string
	if isInterface {
		isIf = "interface "
	}

	return fmt.Sprintf("comparison of pointer to zero-size variable of type %q with %s%q", zType, isIf, withType)
}
