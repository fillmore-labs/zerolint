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
)

type retType struct {
	elem        types.Type
	zeroSized   bool
	valueMethod bool
}

// visitResults examines function bodies for explicit nil return values for zero-sized types.
func (v *Visitor) visitResults(results *ast.FieldList, body *ast.BlockStmt) bool {
	if body == nil { // Skip functions without bodies (e.g., interface methods, external functions)
		return true
	}

	returnTypes, ok := v.hasZeroSizedPointerReturn(results)
	if !ok {
		return true
	}

	ast.Inspect(body, v.inspectBody(returnTypes))

	return true
}

// hasZeroSizedPointerReturn examines a function's return types to detect if any are pointers to zero-sized types.
// It returns a list of return type information and a boolean to indicate if a zero-sized pointer return type exists.
func (v *Visitor) hasZeroSizedPointerReturn(results *ast.FieldList) ([]retType, bool) {
	// Check if the function has return values
	numResults := results.NumFields()
	if numResults == 0 {
		return nil, false
	}

	zeroSizedPointerReturn := make([]retType, numResults)

	var hasZeroSizedPointer bool

	i := 0

	for _, res := range results.List {
		var retInfo retType

		t := v.check.TypesInfo().TypeOf(res.Type)
		retInfo.elem, retInfo.valueMethod, retInfo.zeroSized = v.check.ZeroSizedTypePointer(t)

		if !hasZeroSizedPointer && retInfo.zeroSized {
			hasZeroSizedPointer = true
		}

		names := len(res.Names)
		if names == 0 {
			names = 1
		}

		for range names {
			zeroSizedPointerReturn[i] = retInfo
			i++
		}
	}

	return zeroSizedPointerReturn, hasZeroSizedPointer
}

// inspectBody returns a function suitable for ast.Inspect to walk a function body.
// It looks for `return` statements to check for explicit nil returns for zero-sized pointer types.
func (v *Visitor) inspectBody(returnTypes []retType) func(node ast.Node) bool {
	return func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.FuncLit:
			// Don't check returns in nested function literals
			return false

		case *ast.ReturnStmt:
			v.visitReturnStmt(n, returnTypes)

			return true

		default:
			return true
		}
	}
}

// visitReturnStmt processes `return` statements within a function body.
// It checks if any returned expression corresponding to an expected zero-sized pointer type
// is explicitly 'nil'.
func (v *Visitor) visitReturnStmt(n *ast.ReturnStmt, returnTypes []retType) {
	if len(n.Results) != len(returnTypes) {
		// Skip return statements with differing arity
		return
	}

	for i, result := range n.Results {
		returnType := returnTypes[i]
		if !returnType.zeroSized {
			// Skip if this return position is not expected to be a pointer to a zero-sized type
			continue
		}

		tv, ok := v.check.TypesInfo().Types[result]
		if !ok {
			continue
		}

		switch t := tv.Type.(type) {
		case *types.Basic:
			// Check if the returned expression is the identifier 'nil'
			if t.Kind() == types.UntypedNil {
				cM := msgFormatf(catReturnNil, returnType.valueMethod,
					"explicitly returning nil for pointer to zero-sized type %q", returnType.elem)
				fixes := v.check.ReplaceWithZeroValue(result, returnType.elem)
				v.check.Report(result, cM, fixes)
			}

		case *types.Tuple:
			// Tuples are not handled here. Exiting, since the count is off.
			return
		}
	}
}
