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

package base

import (
	"go/types"
	"strings"
)

// ZeroSizedTypePointer checks whether t is a pointer to a zero-sized type.
// It returns a boolean indicating if t is a pointer to a zero-sized type (and not excluded or an unhandled C type)
// and if true, the underlying element type and another boolean (`valueMethod`) indicating if the zero-sized element
// type has any methods with non-pointer (value) receivers.
func (v *Base) ZeroSizedTypePointer(t types.Type) (elem types.Type, valueMethod, zeroSized bool) {
	if t == nil {
		return nil, false, false
	}

	p, ok := t.Underlying().(*types.Pointer)
	if !ok {
		return nil, false, false // Not a pointer type
	}

	elem = p.Elem()
	valueMethod, zeroSized = v.ZeroSizedType(elem)

	return elem, valueMethod, zeroSized
}

// ZeroSizedType determines whether t is a zero-sized type, considering exclusions and C types.
// It returns true if the type is zero-sized and not excluded.
// It also returns a boolean (`valueMethod`) indicating if a named zero-sized type
// (or a type it embeds that is zero-sized) has methods with non-pointer (value) receivers.
// For unnamed types that don't embed relevant named types, valueMethod is generally false.
// Results for named types are cached to improve performance.
func (v *Base) ZeroSizedType(t types.Type) (valueMethod, zeroSized bool) { //nolint:cyclop
	if t == nil {
		return false, false
	}

	var tn *types.TypeName // The TypeName object corresponding to the type definition (for exclusion and C type checks)
	switch t := t.(type) {
	case *types.Named:
		tn = t.Obj()

	case *types.Alias:
		tn = t.Obj()

	case *types.Basic, *types.Pointer, *types.Interface, *types.Slice, *types.Map, *types.Chan, *types.Signature:
		return false, false // These types are never zero-sized.

	case *types.TypeParam:
		return false, false // Type parameters are unresolved at this stage for size.
	}

	if tn != nil {
		// Heuristic to avoid issues with cgo types like _Ctype_struct_foo
		if strings.HasPrefix(tn.Name(), "_Ctype_") {
			return false, false
		}

		// Check if the type definition is explicitly ignored by its position.
		if v.ignored.Has(tn.Pos()) {
			return false, false
		}
	}

	// Check cache first, [typeutil.Map] resolves aliases
	if c, ok := v.cache.At(t).(typeCache); ok {
		return c.valueMethod, c.zeroSized
	}

	// Recursively check if the underlying type is zero-sized.
	// The `checkMethods` parameter is true here because we are at the top-level call for type `t`.
	vm, zs := isZeroSized(t, true)
	if zs {
		typeString := types.TypeString(t, nil)

		// Track the zero-sized type name for reporting, even if excluded.
		v.Detected.Insert(typeString)

		// Filter out user-specified Excludes by type name.
		if v.Excludes.Has(typeString) {
			zs = false
		}
	}

	// Cache the result.
	v.cache.Set(t, typeCache{zeroSized: zs, valueMethod: vm})

	return vm, zs
}

// isZeroSized determines whether the underlying type of t is provably zero-sized.
// It checks arrays of length 0 and structs where all fields are zero-sized.
// `checkMethods` indicates if method checks should be performed for the current type `t`
// or its embedded fields.
func isZeroSized(t types.Type, checkMethods bool) (valueMethod, zeroSized bool) { //nolint:cyclop
	switch u := t.Underlying().(type) {
	case *types.Array:
		// An array is zero-sized if its length is 0 or its element type is zero-sized.
		// The former (len 0) is the primary case for zero-sized arrays.
		if u.Len() > 0 {
			// We don't propagate `checkMethods`, because it's not relevant for elements of an array.
			_, elemZs := isZeroSized(u.Elem(), false)

			// If length > 0, the array is zero-sized iff its element type is zero-sized.
			if !elemZs {
				return false, false
			}
		}

		var arrayValueMethod bool
		if checkMethods {
			arrayValueMethod = hasValueReceiverMethod(t)
		}

		return arrayValueMethod, true

	case *types.Struct:
		var structValueMethod bool // Accumulates valueMethod from embedded zero-sized fields.

		// A struct is zero-sized if all its fields are zero-sized.
		for i := range u.NumFields() {
			f := u.Field(i)

			// For embedded fields, we need to consider if they contribute value methods.
			// Pass `checkMethods` as true only if we are currently checking methods for the parent struct
			// AND this field is embedded AND we haven't yet found value methods from other fields.
			fieldCheckMethods := checkMethods && f.Embedded() && !structValueMethod

			fieldValMethods, fieldZeroSized := isZeroSized(f.Type(), fieldCheckMethods)
			if !fieldZeroSized {
				return false, false
			}

			if fieldValMethods {
				structValueMethod = true
			}
		}

		// If we are checking methods
		// and haven't found value methods via embedded fields, check methods on the named type itself.
		if checkMethods && !structValueMethod {
			structValueMethod = hasValueReceiverMethod(t)
		}

		// All fields are zero-sized, so the struct is zero-sized.
		return structValueMethod, true

	default:
		// All other types (Basic, Chan, Interface, Map, Pointer, Signature, Slice, TypeParam)
		// are considered not zero-sized by this function or handled by ZeroSizedType directly.
		return false, false
	}
}

// hasValueReceiverMethod checks if a type `t` (which should be a Named type or an alias to one)
// has any methods with a non-pointer (value) receiver.
func hasValueReceiverMethod(t types.Type) bool {
	n, ok := types.Unalias(t).(*types.Named) // Resolve aliases to get the underlying Named type if t is an alias.
	if !ok {
		// Not a named type, so it cannot have methods defined directly on it.
		// If it's an unnamed struct, isZeroSized would check embedded named types.
		return false
	}

	// Iterate over all methods associated with the named type `n`.
	// This does not include methods from embedded types.
	for i := range n.NumMethods() {
		m := n.Method(i)

		recv := m.Signature().Recv()
		if recv == nil { // Should not happen for methods
			continue
		}

		// Check if the receiver is a pointer type.
		// If Underlying() is *types.Pointer, it's a pointer receiver.
		if _, isPtrRecv := recv.Type().Underlying().(*types.Pointer); !isPtrRecv {
			// Found a value (non-pointer) receiver.
			return true
		}
	}

	// Embedded types are checked recursively, so we don't need to check them here.
	return false
}
