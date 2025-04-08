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
	"strings"

	"fillmore-labs.com/zerolint/pkg/internal/checker"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

// errorFunc holds a reference to the Error() method of the standard library error interface.
var errorFunc = types.Universe. //nolint:gochecknoglobals,forcetypeassert
				Lookup("error").Type().Underlying().(*types.Interface).Method(0)

// visitFuncDecl examines method declarations with pointer receivers to zero-sized types.
func (v *Visitor) visitFuncDecl(n *ast.FuncDecl) bool {
	if strings.HasPrefix(n.Name.Name, "_Cfunc_") {
		return false
	}

	v.visitFuncRecv(n)

	if v.level.AtLeast(level.Extended) {
		_ = v.visitResults(n.Type.Results, n.Body)
	}

	return true
}

// visitFuncRecv checks method receivers.
// It reports an issue if a method on a pointer to a zero-sized type is found.
//
// Special handling applies to:
//   - Error() methods: Always checked and reported.
//   - Lock()/Unlock() methods: Star removal is not suggested.
func (v *Visitor) visitFuncRecv(n *ast.FuncDecl) {
	if n.Recv == nil || len(n.Recv.List) != 1 {
		return // Skip non-methods.
	}

	p := n.Recv.List[0].Type
	r := v.check.TypesInfo().TypeOf(p)

	elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(r)
	if !zeroSized {
		return
	}

	var cM checker.CategorizedMessage

	switch {
	case v.isError(n):
		cM = msgFormatf(catError, valueMethod, "error interface implemented on pointer to zero-sized type %q", elem)

	case v.level.Below(level.Extended):
		return

	case v.isLock(n):
		// Lock/Unlock methods on zero-sized types (typedef struct noCopy{}), are common
		// to create an embeddable marker that structures should not be copied.
		// Therefore, don't suggest removing the star from the receiver.
		if s, ok := ast.Unparen(p).(*ast.StarExpr); ok {
			v.check.IgnoreStar(s)
		}

		return

	default:
		// Possible improvement:
		// Do not report implementations of [json.Unmarshaler], [encoding.TextUnmarshaler] or [encoding.BinaryUnmarshaler]
		cM = msgFormatf(catReceiver, valueMethod, "method %s has pointer receiver to zero-sized type %q", n.Name.Name, elem)
	}

	fixes := v.check.RemoveStar(p)
	v.check.Report(p, cM, fixes)
}

// isError determines if the given function declaration represents an implementation
// of the error interface's Error() method by checking both its name and signature.
func (v *Visitor) isError(n *ast.FuncDecl) bool {
	if n.Name.Name != errorFunc.Name() {
		return false
	}

	fn, ok := v.check.TypesInfo().Defs[n.Name].(*types.Func)
	if !ok {
		return false
	}

	return types.Identical(fn.Type(), errorFunc.Type())
}

// isLock determines if the function declaration n is a Lock or Unlock method
// with no parameters and no return values, a common signature for locking mechanisms.
func (v *Visitor) isLock(n *ast.FuncDecl) bool {
	if len(n.Type.Params.List) != 0 || n.Type.Results != nil && len(n.Type.Results.List) != 0 {
		return false
	}

	switch n.Name.Name {
	case "Lock", "Unlock":
		return true

	default:
		return false
	}
}
