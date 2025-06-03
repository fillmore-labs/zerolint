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

package analyzer

import (
	"go/token"
	"go/types"
)

// errorFunc holds a reference to the Error() method of the standard library error interface.
var errorFunc = newErrorFunc() //nolint:gochecknoglobals

// newErrorFunc constructs a new reference to the Error() method.
func newErrorFunc() *types.Func {
	const errorMethodName = "Error"

	var noPkg *types.Package

	results := singleVar(types.Typ[types.String])
	sig := types.NewSignatureType(nil, nil, nil, nil, results, false)

	return types.NewFunc(token.NoPos, noPkg, errorMethodName, sig)
}

// singleVar constructs a tuple with a single unnamed variable of type `t`.
func singleVar(t types.Type) *types.Tuple {
	var noPkg *types.Package

	return types.NewTuple(types.NewVar(token.NoPos, noPkg, "", t))
}
