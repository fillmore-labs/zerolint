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
	"go/types"
)

// zeroSizedTypePointer checks wether t is a pointer to a zero-sized type.
// It returns the element type and true if it is, false otherwise.
func (v Visitor) zeroSizedTypePointer(t types.Type) (types.Type, bool) {
	if p, ok := t.Underlying().(*types.Pointer); ok && v.isZeroSizedType(p.Elem()) {
		return p.Elem(), true
	}

	return nil, false
}

// isZeroSizedType determines whether t is a zero-sized type not excluded from detection.
func (v Visitor) isZeroSizedType(t types.Type) bool {
	if !zeroSized(t) {
		return false
	}

	// zero-sized type, check if the type's name is in the Excludes set.
	name := t.String()
	if v.Excludes.Has(name) {
		return false
	}

	v.Detected.Insert(name)

	return true
}

// zeroSized determines whether t is provably a zero-sized type.
func zeroSized(t types.Type) bool {
	switch x := t.Underlying().(type) {
	case *types.Array:
		// array type, check if the array length is zero or if the element type is zero-sized.
		if x.Len() == 0 {
			return true
		}

		return zeroSized(x.Elem())

	case *types.Struct:
		// struct type, check if all fields have zero-sized types.
		for i := range x.NumFields() {
			if !zeroSized(x.Field(i).Type()) {
				return false
			}
		}

		return true

	/* not really useful and doesn't work with '-fix':
	case *types.Interface:
		// interface type, check if any of the embedded types are zero-sized.
		for i := 0; i < x.NumEmbeddeds(); i++ {
			if zeroSized(x.EmbeddedType(i)) {
				return true
			}
		}

		return false

	case *types.Union:
		// union type, check all variants are zero-sized.
		for i := 0; i < x.Len(); i++ {
			if !zeroSized(x.Term(i).Type()) {
				return false
			}
		}

		return true
	*/

	default:
		return false
	}
}
