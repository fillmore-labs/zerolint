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
)

// operandInfo holds type information for comparison operands.
type operandInfo struct {
	zeroSizedPointer, valueMethod bool
	infoType                      types.Type
}

// visitCmp analyzes comparison expressions (x == y, x != y, errors.Is(x, y)) for comparisons
// involving pointers to zero-sized types.
func (v *Visitor) visitCmp(n ast.Node, x, y ast.Expr) bool {
	var (
		left, right operandInfo
		ok          bool
	)

	if left, ok = v.operandInfo(x); !ok {
		return true
	}

	if right, ok = v.operandInfo(y); !ok {
		return true
	}

	var cM checker.CategorizedMessage

	switch {
	case left.zeroSizedPointer && right.zeroSizedPointer:
		cM = comparisonMessage(left.infoType, right.infoType, left.valueMethod || right.valueMethod)

	case left.zeroSizedPointer:
		cM = comparisonMessagePointerInterface(left.infoType, right.infoType, left.valueMethod)

	case right.zeroSizedPointer:
		cM = comparisonMessagePointerInterface(right.infoType, left.infoType, right.valueMethod)

	default:
		return true
	}

	v.check.Report(n, cM, nil)

	// no fixes, so dive deeper.
	return true
}

// operandInfo extracts type information about comparison operands,
// identifying whether they are pointers to zero-sized types or interfaces.
func (v *Visitor) operandInfo(x ast.Expr) (operandInfo, bool) {
	tv := v.check.TypesInfo().Types[x]
	if tv.IsNil() {
		return operandInfo{}, false // comparisons to nil are not flagged
	}

	t := tv.Type.Underlying()

	if elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(t); zeroSized {
		return operandInfo{zeroSizedPointer: true, valueMethod: valueMethod, infoType: elem},
			true // comparisons with a pointer to zero-size vriable
	}

	if _, ok := t.(*types.Interface); ok {
		return operandInfo{infoType: tv.Type}, true // comparisons with an interface
	}

	return operandInfo{}, false // other comparisons
}

// comparisonMessage generates a diagnostic message for comparing two pointers to zero-sized types.
func comparisonMessage(left, right types.Type, valueMethod bool) checker.CategorizedMessage {
	leftTypeString := types.TypeString(left, nil)
	rightTypeString := types.TypeString(right, nil)

	if leftTypeString == rightTypeString { // types.Identical ignores aliases
		return msgFormatf(catComparison, valueMethod, "comparison of pointers to zero-size type %q", leftTypeString)
	}

	return msgFormatf(catComparison, valueMethod,
		"comparison of pointers to zero-size types %q and %q", leftTypeString, rightTypeString)
}

// comparisonMessagePointerInterface generates a diagnostic message for pointer-to-interface comparison.
func comparisonMessagePointerInterface(elemOp, interfaceOp types.Type, valueMethod bool) checker.CategorizedMessage {
	elemTypeString := types.TypeString(elemOp, nil)
	interfaceTypeString := types.TypeString(interfaceOp, nil)

	if interfaceTypeString == "error" {
		return msgFormatf(catComparisonError, valueMethod,
			"comparison of pointer to zero-size type %q with error interface", elemTypeString)
	}

	return msgFormatf(catComparisonInterface, valueMethod,
		"comparison of pointer to zero-size type %q with interface of type %q", elemTypeString, interfaceTypeString)
}
