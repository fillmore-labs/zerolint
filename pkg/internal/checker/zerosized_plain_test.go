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

package checker_test

import "go/types"

func isZeroSizedPlain(t types.Type) bool {
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
		var top types.Type
		top, stack = stack[len(stack)-1], stack[:len(stack)-1]

		// We could check for excluded *[types.Named] or *[types.Alias] here.

		switch u := top.Underlying().(type) {
		case *types.Array:
			// An array is zero-sized if its length is 0 or its element type is zero-sized.
			if u.Len() > 0 {
				stack = append(stack, u.Elem())
			}

		case *types.Struct:
			// A struct is zero-sized if all its fields are zero-sized.
			for field := range u.Fields() {
				stack = append(stack, field.Type())
			}

		default:
			// Other types (Basic, Chan, Interface, Map, Pointer, Signature, Slice, TypeParam)
			// are not zero-sized.
			return false
		}
	}

	return len(stack) == 0 // All types are zero-sized
}
