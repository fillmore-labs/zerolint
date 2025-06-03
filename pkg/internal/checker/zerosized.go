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

package checker

import "go/types"

// ZeroSizedTypePointer checks whether t is a pointer to a zero-sized type.
//
// It returns:
//   - zeroSized: true if t is a pointer to a zero-sized type (and not excluded).
//
// and, if zeroSized is true:
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
// and, if zeroSized is true:
//   - valueMethod: the type has value receiver methods.
func (c *Checker) ZeroSizedType(t types.Type) (valueMethod, zeroSized bool) {
	if c.isIgnored(t) {
		return false, false
	}

	vM, zS := c.lookupOrCalculate(t)

	// The cached zS value reflects structural zero-sizedness. We can't return it directly,
	// since the specific type name can differ for aliases vs. original types, which gives
	// different results in subsequent filtering based on user-defined exclusions.
	if zS {
		typeName := types.TypeString(t, nil)

		c.trackType(typeName, vM)

		zS = !c.isExcluded(typeName)
	}

	return vM, zS
}

func (c *Checker) lookupOrCalculate(t types.Type) (valueMethod, zeroSized bool) {
	// typeCache holds cached results of zero-sized and value method checks for a type.
	type typeCache struct {
		zeroSized, valueMethod bool
	}

	// Check cache first, [typeutil.Map] resolves aliases
	if cached, ok := c.cache.At(t).(typeCache); ok {
		return cached.valueMethod, cached.zeroSized
	}

	// Check if the underlying type is zero-sized.
	if !IsZeroSized(t) {
		// Cache the result.
		c.cache.Set(t, typeCache{zeroSized: false})

		return false, false
	}
	// Check if the zero-sized type has value receiver methods.
	vM := hasValueMethod(t)

	// Cache the result.
	c.cache.Set(t, typeCache{zeroSized: true, valueMethod: vM})

	return vM, true
}

// Track the zero-sized type name for reporting, even if excluded.
// Excluded types are still detected but not considered for analysis.
func (c *Checker) trackType(typeName string, vM bool) {
	c.Detected[typeName] = vM
}

// Filter out type name by user-specified excludes or regex.
func (c *Checker) isExcluded(typeName string) bool {
	return c.Excludes.Has(typeName) || c.Regex != nil && !c.Regex.MatchString(typeName)
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

// IsZeroSized determines whether the type t is provably zero-sized.
func IsZeroSized(t types.Type) bool {
	const (
		// initialStackCapacity is the initial capacity for the type traversal stack
		// This reduces allocations for most common cases.
		initialStackCapacity = 10

		// maxIterations protects from recusive types (should not happen) and types
		// that are too costly to evaluate.
		maxIterations = 100
	)

	store := [initialStackCapacity]types.Type{t}
	stack := store[:1]

	for budget := maxIterations; budget > 0 && len(stack) > 0; budget-- {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch u := top.Underlying().(type) {
		case *types.Array:
			// An array is zero-sized if its length is 0 or its element type is zero-sized.
			if u.Len() > 0 {
				stack = append(stack, u.Elem())
			}

		case *types.Struct:
			// A struct is zero-sized if all its fields are zero-sized.
			for i := range u.NumFields() {
				stack = append(stack, u.Field(i).Type())
			}

		default:
			// Other types (Basic, Chan, Interface, Map, Pointer, Signature, Slice, TypeParam)
			// are not zero-sized.
			return false
		}
	}

	return len(stack) == 0 // All types are zero-sized
}

// hasValueMethod checks if a type has any methods with a value receiver.
func hasValueMethod(t types.Type) bool {
	mset := types.NewMethodSet(t)

	return mset.Len() > 0
}
