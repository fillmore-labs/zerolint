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
	"strings"

	"golang.org/x/tools/go/ast/inspector"

	"fillmore-labs.com/zerolint/pkg/internal/analyzer/msg"
	"fillmore-labs.com/zerolint/pkg/internal/diag"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

// visitFuncDecl examines method declarations with pointer receivers to zero-sized types.
func (v *visitor) visitFuncDecl(c inspector.Cursor, n *ast.FuncDecl) bool {
	if isCfunc(n.Name.Name) {
		return false
	}

	v.visitFuncRecv(n)

	if v.level.AtLeast(level.Extended) {
		v.checkReturns(c, n.Body, n.Type.Results)
	}

	return true
}

// isCtype is a heuristic to avoid issues with cgo functions like _Cfunc_foo.
func isCfunc(name string) bool {
	const cgoFuncPrefix = "_Cfunc_"

	return strings.HasPrefix(name, cgoFuncPrefix)
}

// visitFuncRecv checks method receivers.
// It reports an issue if a method on a pointer to a zero-sized type is found.
//
// Special handling applies to:
//   - Error() methods: Always checked and reported.
//   - Lock()/Unlock() methods: Star removal is not suggested.
func (v *visitor) visitFuncRecv(n *ast.FuncDecl) {
	if n.Recv == nil || len(n.Recv.List) != 1 {
		return // Skip non-methods.
	}

	p := n.Recv.List[0].Type
	r := v.diag.TypesInfo().TypeOf(p)

	elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(r)
	if !zeroSized {
		return
	}

	var cM diag.CategorizedMessage

	switch {
	case v.isError(n):
		cM = msg.Formatf(msg.CatError, valueMethod, "error interface implemented on pointer to zero-sized type %q", elem)

	case v.level.Below(level.Extended):
		return

	case isLock(n, elem):
		// Lock/Unlock methods on zero-sized types (typedef struct noCopy{}), are common
		// to create an embeddable marker that structures should not be copied.
		// Therefore, don't suggest removing the star from the receiver.
		if s, ok := ast.Unparen(p).(*ast.StarExpr); ok {
			v.ignoreStar(s)
		}

		return

	default:
		// Possible improvement:
		// Do not report implementations of [json.Unmarshaler], [encoding.TextUnmarshaler] or [encoding.BinaryUnmarshaler]
		cM = msg.Formatf(msg.CatReceiver, valueMethod,
			"method %s has pointer receiver to zero-sized type %q", n.Name.Name, elem)
	}

	fixes := v.removeStar(p)
	v.diag.Report(p, cM, fixes)
}

// isError determines if the given function declaration represents an implementation
// of the error interface's Error() method by checking both its name and signature.
func (v *visitor) isError(n *ast.FuncDecl) bool {
	if n.Name.Name != errorFunc.Name() {
		return false
	}

	def := v.diag.TypesInfo().Defs[n.Name]
	fn, ok := def.(*types.Func)

	if !ok { // should not happen
		v.diag.LogErrorf(n, "expected *types.Func, got %T", def)

		return false
	}

	return types.Identical(fn.Type(), errorFunc.Type())
}

// isLock determines if the function declaration represents a Lock or Unlock
// method on a pointer receiver to a struct type. This pattern is often used with
// a zero-sized `noCopy` struct for embedding, and the receiver type is a pointer to that struct.
func isLock(n *ast.FuncDecl, elem types.Type) bool {
	if hasParams(n.Type) || hasResults(n.Type) {
		return false // A method with parameters or results is not a standard Lock/Unlock.
	}

	switch n.Name.Name {
	case "Lock", "Unlock":
		_, structPtrReceiver := elem.Underlying().(*types.Struct)

		return structPtrReceiver // Only valid on struct types.

	default:
		return false
	}
}

// hasParams checks if the given function type has any parameters.
func hasParams(f *ast.FuncType) bool {
	return len(f.Params.List) != 0
}

// hasResults checks if the given function type has any results.
func hasResults(f *ast.FuncType) bool {
	return f.Results != nil && len(f.Results.List) != 0
}
