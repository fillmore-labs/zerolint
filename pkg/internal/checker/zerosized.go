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

package checker

import "go/types"

// ZeroSizedTypePointer checks whether t is a pointer to a zero-sized type.
//
// It returns:
//   - zeroSized: true if t is a pointer to a zero-sized type (and not excluded).
//
// and if zeroSized is true:
//   - elem: the underlying zero-sized element type of the pointer.
//   - valueMethod: the element type has value receiver methods.
func (c *Checker) ZeroSizedTypePointer(t types.Type) (elem types.Type, valueMethod, zeroSized bool) {
	if t == nil {
		return nil, false, false
	}

	p, ok := t.Underlying().(*types.Pointer)
	if !ok {
		return nil, false, false // Not a pointer type
	}

	elem = p.Elem()
	valueMethod, zeroSized = c.ZeroSizedType(elem)

	return elem, valueMethod, zeroSized
}

// ZeroSizedType determines whether t is a zero-sized type, considering exclusions.
//
// It returns:
//   - zeroSized: true if t is zero-sized (and not excluded).
//
// and if zeroSized is true:
//   - valueMethod: the type has value receiver methods.
func (c *Checker) ZeroSizedType(t types.Type) (valueMethod, zeroSized bool) {
	if c.isIgnored(t) {
		return false, false
	}

	// typeCache holds cached results of zero-sized and value method checks for a type.
	type typeCache struct {
		zeroSized, valueMethod bool
	}

	var zS, vM bool
	// Check cache first, [typeutil.Map] resolves aliases
	if cached, ok := c.cache.At(t).(typeCache); ok {
		zS, vM = cached.zeroSized, cached.valueMethod
	} else {
		// Recursively check if the underlying type is zero-sized.
		zS = isZeroSized(t)
		if zS { // Check if the zero-sized type has value receiver methods.
			vM = hasValueMethod(t)
		}

		// Cache the result.
		c.cache.Set(t, typeCache{zeroSized: zS, valueMethod: vM})
	}

	// The cached zS value reflects structural zero-sizedness. We can't return it directly,
	// since the specific type name can differ for aliases vs. original types, which gives
	// different results in subsequent filtering based on user-defined exclusions.
	if zS {
		typeString := types.TypeString(t, nil)

		detectedType := typeString
		if vM {
			detectedType += " (value methods)"
		}

		// Track the zero-sized type name for reporting, even if excluded.
		// Excluded types are still detected but not considered for analysis.
		c.Detected.Insert(detectedType)

		// Filter out type name by user-specified excludes or regex.
		if c.Excludes.Has(typeString) || c.Regex != nil && !c.Regex.MatchString(typeString) {
			zS = false // Excluded by the user-specified list or regex.
		}
	}

	return vM, zS
}

// isIgnored checks if a type should be ignored by the zero-size analysis
// (e.g., explicitly excluded via `//nolint:zerolint` directive or not a candidate type).
func (c *Checker) isIgnored(t types.Type) bool {
	if t == nil {
		return true
	}

	var tn *types.TypeName
	switch t := t.(type) {
	case *types.Named:
		tn = t.Obj()
		// Check below

	case *types.Alias:
		tn = t.Obj()
		// Check below

	case *types.Struct, *types.Array:
		return false

	default:
		// Other types (Basic, Chan, Interface, Map, Pointer, Signature, Slice, TypeParam)
		// are not zero-sized, so they are not considered for analysis.
		return true
	}

	// Check if the type definition is explicitly excluded.
	return c.ExcludedTypeDefs.ExcludedType(tn)
}

// isZeroSized determines recursively whether the underlying type of t is provably zero-sized.
// It checks arrays of length 0 and structs where all fields are zero-sized.
func isZeroSized(t types.Type) bool {
	if t == nil {
		return false
	}

	switch u := t.Underlying().(type) {
	case *types.Array:
		// An array is zero-sized if its length is 0 or its element type is zero-sized.
		// The former (len 0) is the primary case for zero-sized arrays.
		return u.Len() == 0 || isZeroSized(u.Elem())

	case *types.Struct:
		// A struct is zero-sized if all its fields are zero-sized.
		for i := range u.NumFields() {
			f := u.Field(i)
			if !isZeroSized(f.Type()) {
				return false
			}
		}

		// All fields are zero-sized, so the struct is zero-sized.
		return true

	default:
		// Other types (Basic, Chan, Interface, Map, Pointer, Signature, Slice, TypeParam)
		// are not zero-sized.
		return false
	}
}

// hasValueMethod checks if a type has any methods with a value receiver.
func hasValueMethod(t types.Type) bool {
	mset := types.NewMethodSet(t)

	return mset.Len() > 0
}
