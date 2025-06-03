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
)

// visitCallSelection handles method expressions selected from a value (dec.Decode(...)) or type
// ((*json.Decoder).Decode).
// For method values, it delegates to visitMethodVal.
// For method expressions with receivers that are pointers to zero-sized types, it reports an issue.
func (v *visitor) visitCallSelection(fun *ast.SelectorExpr, sel *types.Selection) bool {
	switch sel.Kind() {
	case types.MethodVal:
		f, ok := sel.Obj().(*types.Func)
		if !ok { // should not happen
			v.diag.LogErrorf(fun, "Expected *types.Func, got %T", sel.Obj())

			return true
		}

		// Delegate selections like encoding/json#Decoder.Decode.
		return visitMethodVal(f)

	case types.MethodExpr:
		// Method used as a function value (e.g., (*T).Method).
		recv := sel.Recv()
		// We care about the receiver of the method expression, *not* the signature of the possibly embedded function:
		//	if fun, ok := sel.Obj().(*types.Func); ok {
		//		recv = fun.Signature().Recv().Type()
		//	}

		if isJSONDecoder(recv) {
			return false // Do not report pointers in json.Decoder#Decode.
		}

		if elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(recv); zeroSized {
			cM := msg.Formatf(msg.CatMethodExpression, valueMethod,
				"method expression receiver is pointer to zero-size type %q", elem)
			fixes := v.removeStar(fun.X)
			v.diag.Report(fun, cM, fixes)
		}

		return true

	case types.FieldVal:
		// struct field with function type
	}

	return true
}

// visitMethodVal checks for method calls on specific receivers, particularly
// looking for the Decode method on json.Decoder, which requires pointers.
func visitMethodVal(fun *types.Func) bool {
	if fun.Name() != "Decode" { // We are only interested in `Decode`.
		return true
	}

	// Use the method's declared receiver type rather than the selection's receiver
	// type (sel.Recv()) to correctly identify [json.Decoder], even if it's embedded.
	recv := fun.Signature().Recv().Type()

	// Do not report pointers in json.Decoder#Decode.
	return !isJSONDecoder(recv)
}

func isJSONDecoder(recv types.Type) bool {
	// extract the underlying named type from a pointer type.
	ptr, ok := recv.Underlying().(*types.Pointer)
	if !ok {
		return false
	}

	named, ok := types.Unalias(ptr.Elem()).(*types.Named)
	if !ok { // Since we are only calling for method selectors, this should be a named type.
		return false
	}

	tn := named.Obj()

	// Check for method receiver *encoding/json.Decoder.
	return tn.Pkg().Path() == "encoding/json" && tn.Name() == "Decoder"
}
