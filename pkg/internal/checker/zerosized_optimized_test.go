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

// isZeroSizedSemiOptimized avoids push/pop in simple cases.
func isZeroSizedSemiOptimized(t types.Type) bool {
	const (
		// initialStackCapacity is the initial capacity for the type traversal stack
		// This reduces allocations for most common cases.
		initialStackCapacity = 10

		// maxIterations protects from recusive types (should not happen) and types
		// that are too costly to evaluate.
		maxIterations = 100
	)

	stack := make([]types.Type, 0, initialStackCapacity)

	top := t
	for range maxIterations {
		switch u := top.Underlying().(type) {
		case *types.Array:
			// An array is zero-sized if its length is 0 or its element type is zero-sized.
			if u.Len() > 0 {
				top = u.Elem()

				continue
			}

		case *types.Struct:
			// A struct is zero-sized if all its fields are zero-sized.
			if nf := u.NumFields(); nf > 0 {
				top = u.Field(0).Type()

				for i := 1; i < nf; i++ {
					stack = append(stack, u.Field(i).Type())
				}

				continue
			}

		default:
			// Other types (Basic, Chan, Interface, Map, Pointer, Signature, Slice, TypeParam)
			// are not zero-sized.
			return false
		}

		l := len(stack)
		if l == 0 {
			return true // all types are zero-sized
		}

		// pop next type to check
		top = stack[l-1]
		stack = stack[:l-1]
	}

	return false // too expensive
}

// isZeroSizedOptimized uses a stack only for deeply nested structs.
func isZeroSizedOptimized(t types.Type) bool {
	for typ := t; ; {
		switch u := typ.Underlying().(type) {
		case *types.Array:
			if u.Len() == 0 {
				return true
			}

			typ = u.Elem()

		case *types.Struct:
			switch u.NumFields() {
			case 0:
				return true

			case 1:
				typ = u.Field(0).Type()

			default:
				return isZeroSizedStructOnly(u)
			}

		default:
			return false
		}
	}
}

func isZeroSizedStructOnly(s *types.Struct) bool {
	const initialStackCapacity = 10

	for top, stack := s, make([]*types.Struct, 0, initialStackCapacity); ; {
		for field := range top.Fields() {
		fieldLoop:
			for ft := field.Type(); ; {
				switch uft := ft.Underlying().(type) {
				case *types.Array:
					if uft.Len() == 0 {
						break fieldLoop
					}

					ft = uft.Elem()

				case *types.Struct:
					stack = append(stack, uft)

					break fieldLoop

				default:
					return false
				}
			}
		}

		l := len(stack)
		if l == 0 {
			return true
		}

		top = stack[l-1]
		stack = stack[:l-1]
	}
}
