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

// ZeroMsg returns a message for unnamed function parameters pointing to zero-sized types.
func (msgParam) ZeroMsg(t types.Type) (category, string) {
	return catParameter, fmt.Sprintf("function has pointer parameter to zero-sized type %q", t)
}

// SingularMsg returns a message for a single named function parameter pointing to a zero-sized type.
func (msgParam) SingularMsg(name string, t types.Type) (category, string) {
	return catParameter, fmt.Sprintf("function parameter %q points to zero-sized type %q", name, t)
}

// PluralMsg returns a message for multiple function parameters pointing to a zero-sized type.
func (msgParam) PluralMsg(names string, t types.Type) (category, string) {
	return catParameter, fmt.Sprintf("function parameters %s point to zero-sized type %q", names, t)
}

// msgResult implements msgFormatter for function results, providing appropriate diagnostic
// messages for results that are pointers to zero-sized types.
type msgResult struct{}

// ZeroMsg returns a message for unnamed function results pointing to zero-sized types.
func (msgResult) ZeroMsg(t types.Type) (category, string) {
	return catResult, fmt.Sprintf("function has pointer result to zero-sized type %q", t)
}

// SingularMsg returns a message for a single named function result pointing to a zero-sized type.
func (msgResult) SingularMsg(name string, t types.Type) (category, string) {
	return catResult, fmt.Sprintf("function result %q points to zero-sized type %q", name, t)
}

// PluralMsg returns a message for multiple function results pointing to a zero-sized type.
func (msgResult) PluralMsg(names string, t types.Type) (category, string) {
	return catResult, fmt.Sprintf("function results %s point to zero-sized type %q", names, t)
}
