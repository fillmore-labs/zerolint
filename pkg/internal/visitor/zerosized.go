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
	"iter"
)

// zeroSizedTypePointer checks whether t is a pointer to a zero-sized type.
// It returns the element's type and true if it is zero-sized, false otherwise.
func (v *Visitor) zeroSizedTypePointer(t types.Type) (types.Type, bool) {
	if t == nil {
		return nil, false
	}
	if p, ok := t.Underlying().(*types.Pointer); ok && v.zeroSizedType(p.Elem()) {
		return p.Elem(), true
	}

	return nil, false
}

// zeroSizedType determines whether t is a zero-sized type not excluded from detection.
func (v *Visitor) zeroSizedType(t types.Type) bool {
	if t == nil || !zeroSized(t) {
		return false
	}

	if n, ok := t.(*types.Named); ok {
		// Check if declared in a generated file or excluded by comment
		if tn := n.Obj(); v.ignored.Has(tn.Pos()) {
			return false
		}
	}

	name := t.String()
	v.detected.Insert(name)

	// zero-sized type, check if the type's name is in the Excludes set.
	return !v.excludes.Has(name)
}

// zeroSized determines whether t is provably a zero-sized type.
func zeroSized(t types.Type) bool {
	switch x := t.Underlying().(type) {
	case *types.Array:
		// array type, check if the array length is zero or if the element type is zero-sized.
		return x.Len() == 0 || zeroSized(x.Elem())

	case *types.Struct:
		// struct type, check if all fields have zero-sized types.
		for f := range allFields(x) {
			if ft := f.Type(); ft == nil || !zeroSized(ft) {
				return false
			}
		}

		return true

	default:
		return false
	}
}

// allFields returns an iterator over all fields of the struct.
func allFields(s *types.Struct) iter.Seq[*types.Var] {
	return func(yield func(*types.Var) bool) {
		for i := range s.NumFields() {
			if !yield(s.Field(i)) {
				break
			}
		}
	}
}
