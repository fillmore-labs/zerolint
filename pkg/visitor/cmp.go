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

func (v Visitor) visitCmp(n ast.Node, x, y ast.Expr) bool { //nolint:cyclop
	var p [2]struct {
		name          string
		isZeroPointer bool
		isInterface   bool
	}

	for i, z := range []ast.Expr{x, y} {
		t := v.TypesInfo.Types[z]
		if t.IsNil() {
			return true
		}
		switch x := t.Type.Underlying().(type) {
		case *types.Pointer:
			e := x.Elem()
			p[i].name = e.String()
			p[i].isZeroPointer = v.isZeroSizeType(e)

		case *types.Interface:
			p[i].name = t.Type.String()
			p[i].isInterface = true

		default:
			p[i].name = t.Type.String()
		}
	}

	var message string
	switch {
	case p[0].isZeroPointer && p[1].isZeroPointer:
		message = comparisonMessage(p[0].name, p[1].name)

	case p[0].isZeroPointer:
		message = comparisonIMessage(p[0].name, p[1].name, p[1].isInterface)

	case p[1].isZeroPointer:
		message = comparisonIMessage(p[1].name, p[0].name, p[0].isInterface)

	default:
		return true
	}

	v.report(n, message, nil)

	return true
}

func comparisonMessage(xType, yType string) string {
	if xType == yType {
		return fmt.Sprintf("comparison of pointers to zero-size variables of type %q", xType)
	}

	return fmt.Sprintf("comparison of pointers to zero-size variables of types %q and %q", xType, yType)
}

func comparisonIMessage(zType, withType string, isInterface bool) string {
	var isIf string
	if isInterface {
		isIf = "interface "
	}

	return fmt.Sprintf("comparison of pointer to zero-size variable of type %q with %s%q", zType, isIf, withType)
}
