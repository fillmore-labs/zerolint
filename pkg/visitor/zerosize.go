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
	"go/ast"
	"go/types"
	"log"
)

func (v Visitor) isZeroPointer(x ast.Expr) bool {
	t := v.TypesInfo.Types[x].Type

	return v.isZeroPointerType(t)
}

func (v Visitor) isZeroPointerType(t types.Type) bool {
	p, ok := t.(*types.Pointer)
	if !ok {
		return false
	}

	return v.isZeroSizeType(p.Elem())
}

func (v Visitor) isZeroSizeType(t types.Type) bool {
	if !zeroSize(t) {
		return false
	}

	name := t.String()
	if _, ok := v.Excludes[name]; ok {
		return false
	}

	if v.ZeroTrace {
		log.Printf("found zero-size type %q", t.String())
	}

	return true
}

func zeroSize(t types.Type) bool {
	switch x := t.Underlying().(type) {
	case *types.Array:
		if x.Len() == 0 {
			return true
		}

		return zeroSize(x.Elem())

	case *types.Struct:
		n := x.NumFields()
		for i := 0; i < n; i++ {
			if !zeroSize(x.Field(i).Type()) {
				return false
			}
		}

		return true

	default:
		return false
	}
}
