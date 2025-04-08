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

// errorFunc holds a reference to the Error() method of the standard library error interface.
var errorFunc = types.Universe. //nolint:gochecknoglobals,forcetypeassert
				Lookup("error").Type().Underlying().(*types.Interface).Method(0)

// visitFuncDecl examines function declarations for methods with pointer receivers to zero-sized types.
func (v *Visitor) visitFuncDecl(n *ast.FuncDecl) bool {
	if n.Recv == nil || len(n.Recv.List) != 1 {
		return true
	}

	p := n.Recv.List[0].Type
	r := v.Pass.TypesInfo.TypeOf(n.Recv.List[0].Type)
	t, ok := v.zeroSizedTypePointer(r)
	if !ok {
		return true
	}

	var message string
	if !v.isError(n.Name) {
		if !v.full && !v.method {
			return true
		}
		message = fmt.Sprintf("method %s has pointer receiver to zero-sized type %q (ZL04)", n.Name.Name, t)
	} else {
		message = fmt.Sprintf("error interface implemented on pointer to zero-sized type %q (ZL03)", t)
	}
	fixes := v.removeStar(p)
	v.report(p, message, fixes)

	return true
}

// isError determines if the given method name identifier represents an implementation
// of the error interface's Error() method, checking both the name and signature.
func (v *Visitor) isError(name *ast.Ident) bool {
	if name.Name != errorFunc.Name() {
		return false
	}
	def := v.Pass.TypesInfo.Defs[name]
	if def == nil {
		return true
	}

	return types.Identical(def.Type(), errorFunc.Type())
}
