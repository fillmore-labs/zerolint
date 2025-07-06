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

	"fillmore-labs.com/zerolint/pkg/internal/analyzer/msg"
	"fillmore-labs.com/zerolint/pkg/internal/diag"
)

// visitCallArgs checks for explicit nil arguments to pointers to zero-sized parameters.
func (v *Visitor) visitCallArgs(sig *types.Signature, args []ast.Expr) {
	// Matching parameter for each argument. Frozen at the last parameter for variadic functions.
	var (
		param       *types.Var
		elem        types.Type
		valueMethod bool
		zeroSized   bool
	)

	params := sig.Params()
	variadic := sig.Variadic()

	for i, arg := range args {
		check := true

		switch {
		case variadic && i == params.Len()-1: // Variadic and last parameter: Freeze at slice element.
			param = params.At(i)
			if variadicType, ok := param.Type().(*types.Slice); ok {
				elem, valueMethod, zeroSized = v.Check.ZeroSizedTypePointer(variadicType.Elem())
			} else { // variadic arguments should be slices (or string for append, but we are not called for builtins)
				check = false
			}

		case i < params.Len(): // Normal parameters
			param = params.At(i)
			paramType := param.Type()
			elem, valueMethod, zeroSized = v.Check.ZeroSizedTypePointer(paramType)

		case !variadic: // i >= params.Len(): Do nothing for variadic function, otherwise abort (should not happen)
			check = false
		}

		if !check { // should not happen
			v.Diag.LogErrorf(arg,
				"Argument count mismatch, argument %d of %d parameters (variadic: %t)", i, params.Len(), variadic)

			break
		}

		if !zeroSized { // parameter not a pointer to a zero-sized type
			continue
		}

		if tv, ok := v.Diag.TypesInfo().Types[arg]; ok && tv.IsNil() {
			v.handleNilArg(arg, param.Name(), elem, valueMethod)
		}
	}
}

func (v *Visitor) handleNilArg(arg ast.Expr, name string, elem types.Type, valueMethod bool) {
	var cM diag.CategorizedMessage
	if name == "" {
		cM = msg.Formatf(msg.CatArgumentNil, valueMethod, "passing explicit nil as parameter of type %q", elem)
	} else {
		cM = msg.Formatf(msg.CatArgumentNil, valueMethod,
			"passing explicit nil as parameter %q pointing to zero-sized type %q", name, elem)
	}

	fixes := v.Diag.ReplaceWithZeroValue(arg, elem)
	v.Diag.Report(arg, cM, fixes)
}
