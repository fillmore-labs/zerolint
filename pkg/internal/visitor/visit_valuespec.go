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
	"strings"

	"fillmore-labs.com/zerolint/pkg/internal/checker"
)

// visitValueSpec analyzes variable declarations (`var` or `const` specs)
// to detect if they explicitly declare variables as pointers to zero-sized types.
func (v *Visitor) visitValueSpec(n *ast.ValueSpec) bool {
	if n.Type == nil {
		return true
	}

	t, ok := v.check.TypesInfo().Types[n.Type]
	if !ok {
		return true
	}

	elem, valueMethod, zeroSized := v.check.ZeroSizedTypePointer(t.Type)
	if !zeroSized {
		return true
	}

	for _, name := range n.Names {
		if strings.HasPrefix(name.Name, "_cgo") {
			return false
		}
	}

	formatter := msgValue{}
	cM := formatMessage(formatter, elem, valueMethod, n.Names)
	fixes := v.check.RemoveStar(n.Type) // Suggest removing the pointer '*' from the type declaration
	v.check.Report(n, cM, fixes)

	return true
}

type msgValue struct{}

// zeroMsg returns a message for an unnamed variable pointing to a zero-sized type. This is not possible in Go.
func (msgValue) zeroMsg(typ types.Type, valueMethod bool) checker.CategorizedMessage {
	return msgFormatf(catVar, valueMethod, "variable is pointer to zero-sized type %q", typ)
}

// singularMsg returns a message for a single named variable pointing to a zero-sized type.
func (msgValue) singularMsg(typ types.Type, valueMethod bool, name string) checker.CategorizedMessage {
	return msgFormatf(catVar, valueMethod, "variable %q is pointer to zero-sized type %q", name, typ)
}

// pluralMsg returns a message for multiple variables pointing to a zero-sized type.
func (msgValue) pluralMsg(typ types.Type, valueMethod bool, names string) checker.CategorizedMessage {
	return msgFormatf(catVar, valueMethod, "variables %s point to zero-sized type %q", names, typ)
}
