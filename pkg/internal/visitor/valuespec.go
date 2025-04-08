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
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"fillmore-labs.com/zerolint/pkg/internal/base"
)

// visitValueSpec analyzes variable declarations (`var` or `const` specs)
// to detect if they explicitly declare variables as pointers to zero-sized types.
func (v *Visitor) visitValueSpec(n *ast.ValueSpec) bool {
	if n.Type == nil {
		return true
	}

	t, ok := v.base.TypesInfo().Types[n.Type]
	if !ok {
		return true
	}

	elem, valueMethod, zeroSized := v.base.ZeroSizedTypePointer(t.Type)
	if !zeroSized {
		return true
	}

	for _, name := range n.Names {
		if strings.HasPrefix(name.Name, "_cgo") {
			return false
		}
	}

	cat, message := base.FormatMessage(n.Names, msgValue{}, elem)
	// Suggest removing the pointer '*' from the type declaration
	fixes := v.base.RemoveStar(n.Type)
	v.base.Report(n, cat, valueMethod, message, fixes)

	return true
}

type msgValue struct{}

// ZeroMsg returns a message for an unnamed variable pointing to a zero-sized type.
func (msgValue) ZeroMsg(t types.Type) (category, string) {
	return catVar, fmt.Sprintf("variable is pointer to zero-sized type %q", t)
}

// SingularMsg returns a message for a single named variable pointing to a zero-sized type.
func (msgValue) SingularMsg(name string, t types.Type) (category, string) {
	return catVar, fmt.Sprintf("variable %q is pointer to zero-sized type %q", name, t)
}

// PluralMsg returns a message for multiple variables pointing to a zero-sized type.
func (msgValue) PluralMsg(names string, t types.Type) (category, string) {
	return catVar, fmt.Sprintf("variables %s point to zero-sized type %q", names, t)
}
