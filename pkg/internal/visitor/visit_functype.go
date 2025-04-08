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

	"fillmore-labs.com/zerolint/pkg/internal/checker"
)

// visitFuncType examines function types for parameters and results that are pointers to zero-sized types.
func (v *Visitor) visitFuncType(n *ast.FuncType) bool {
	v.checkFieldList(n.Params, false, msgParam{})

	if n.Results != nil {
		v.checkFieldList(n.Results, false, msgResult{})
	}

	return true
}

// msgParam implements msgFormatter for function parameters, providing appropriate diagnostic
// messages for parameters that are pointers to zero-sized types.
type msgParam struct{}

// zeroMsg returns a message for unnamed function parameters pointing to zero-sized types.
func (msgParam) zeroMsg(typ types.Type, valueMethod bool) checker.CategorizedMessage {
	return msgFormatf(catParameter, valueMethod, "function has pointer parameter to zero-sized type %q", typ)
}

// singularMsg returns a message for a single named function parameter pointing to a zero-sized type.
func (msgParam) singularMsg(typ types.Type, valueMethod bool, name string) checker.CategorizedMessage {
	return msgFormatf(catParameter, valueMethod, "function parameter %q points to zero-sized type %q", name, typ)
}

// pluralMsg returns a message for multiple function parameters pointing to a zero-sized type.
func (msgParam) pluralMsg(typ types.Type, valueMethod bool, names string) checker.CategorizedMessage {
	return msgFormatf(catParameter, valueMethod, "function parameters %s point to zero-sized type %q", names, typ)
}

// msgResult implements msgFormatter for function results, providing appropriate diagnostic
// messages for results that are pointers to zero-sized types.
type msgResult struct{}

// zeroMsg returns a message for unnamed function results pointing to zero-sized types.
func (msgResult) zeroMsg(typ types.Type, valueMethod bool) checker.CategorizedMessage {
	return msgFormatf(catResult, valueMethod, "function has pointer result to zero-sized type %q", typ)
}

// singularMsg returns a message for a single named function result pointing to a zero-sized type.
func (msgResult) singularMsg(typ types.Type, valueMethod bool, name string) checker.CategorizedMessage {
	return msgFormatf(catResult, valueMethod, "function result %q points to zero-sized type %q", name, typ)
}

// pluralMsg returns a message for multiple function results pointing to a zero-sized type.
func (msgResult) pluralMsg(typ types.Type, valueMethod bool, names string) checker.CategorizedMessage {
	return msgFormatf(catResult, valueMethod, "function results %s point to zero-sized type %q", names, typ)
}
