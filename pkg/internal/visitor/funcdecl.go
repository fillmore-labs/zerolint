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
	"strings"

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

func (v *Visitor) visitFuncRecv(n *ast.FuncDecl) {
	if n.Recv == nil || len(n.Recv.List) != 1 {
		return
	}

	p := n.Recv.List[0].Type
	r := v.check.TypesInfo().TypeOf(p)

	elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(r)
	if !zeroSized {
		return
	}

	var (
		cat     category
		message string
	)

	switch {
	case v.isError(n.Name):
		cat = catError
		message = fmt.Sprintf("error interface implemented on pointer to zero-sized type %q", elem)

	case v.level.AtLeast(level.Extended):
		// Possible improvement:
		// Do not report implementations of [json.Unmarshaler], [encoding.TextUnmarshaler] or [encoding.BinaryUnmarshaler]
		cat = catReceiver
		message = fmt.Sprintf("method %s has pointer receiver to zero-sized type %q", n.Name.Name, elem)

	default:
		return
	}

	fixes := v.check.RemoveStar(p)
	v.check.Report(p, cat, valueMethod, message, fixes)
}

// isError determines if the given method name identifier represents an implementation
// of the error interface's Error() method, checking both the name and signature.
func (v *Visitor) isError(name *ast.Ident) bool {
	if name.Name != errorFunc.Name() {
		return false
	}

	def := v.check.TypesInfo().Defs[name]
	if def == nil {
		return true
	}

	return types.Identical(def.Type(), errorFunc.Type())
}
