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
	"go/types"

	"golang.org/x/tools/go/ast/inspector"

	"fillmore-labs.com/zerolint/pkg/internal/analyzer/msg"
)

type retType struct {
	elem        types.Type
	zeroSized   bool
	valueMethod bool
}

// checkReturns examines function bodies for explicit nil return values for pointers to zero-sized types.
func (v *visitor) checkReturns(c inspector.Cursor, body *ast.BlockStmt, results *ast.FieldList) {
	if body == nil {
		return // Skip functions without bodies (e.g. external functions)
	}

	returnTypes, ok := v.hasZeroSizedPointerReturn(results)
	if !ok {
		return // No pointer to zero-sized types
	}

	b, ok := c.FindNode(body)
	if !ok { // should not happen
		v.diag.LogErrorf(body, "Can't find body")

		return
	}

	// inspectBody is a function suitable for [inspector.Cursor] to walk a function body.
	nodes := []ast.Node{(*ast.FuncLit)(nil), (*ast.ReturnStmt)(nil)}
	inspectBody := func(c inspector.Cursor) (descend bool) {
		switch n := c.Node().(type) {
		case *ast.FuncLit:
			// Don't check returns in nested function literals
			return false

		case *ast.ReturnStmt:
			// Check for explicit nil returns for zero-sized pointer types.
			v.visitReturnStmt(n, returnTypes)
		}

		return true
	}

	b.Inspect(nodes, inspectBody)
}

// hasZeroSizedPointerReturn examines a function's return types to detect if any are pointers to zero-sized types.
// It returns a list of return type information and a boolean to indicate if a zero-sized pointer return type exists.
func (v *visitor) hasZeroSizedPointerReturn(results *ast.FieldList) ([]retType, bool) {
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

		t := v.diag.TypesInfo().TypeOf(res.Type)
		retInfo.elem, retInfo.valueMethod, retInfo.zeroSized = v.check.ZeroSizedTypePointer(t)

		if !hasZeroSizedPointer && retInfo.zeroSized {
			hasZeroSizedPointer = true
		}

		retVals := len(res.Names)
		if retVals == 0 { // unnamed result
			retVals = 1
		}

		for range retVals {
			zeroSizedPointerReturn[i] = retInfo
			i++
		}
	}

	return zeroSizedPointerReturn, hasZeroSizedPointer
}

// visitReturnStmt processes `return` statements within a function body.
// It checks if any returned expression corresponding to an expected zero-sized pointer type
// is explicitly 'nil'.
func (v *visitor) visitReturnStmt(n *ast.ReturnStmt, returnTypes []retType) {
	if len(n.Results) != len(returnTypes) {
		// Skip return statements with differing arity
		return
	}

	for i, result := range n.Results {
		returnType := returnTypes[i]
		if !returnType.zeroSized {
			continue // Skip if this return position is not a pointer to a zero-sized type
		}

		// Check if the returned expression is the identifier 'nil'
		tv, ok := v.diag.TypesInfo().Types[result]
		if !ok { // should not happen
			v.diag.LogErrorf(result, "Unknown result type")

			continue
		}

		if tv.IsNil() {
			cM := msg.Formatf(msg.CatReturnNil, returnType.valueMethod,
				"explicitly returning nil for pointer to zero-sized type %q", returnType.elem)
			fixes := v.diag.ReplaceWithZeroValue(result, returnType.elem)
			v.diag.Report(result, cM, fixes)
		}
	}
}
