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
	"go/token"
)

// visitUnary checks expressions in form &x.
func (v Visitor) visitUnary(x *ast.UnaryExpr) bool {
	if x.Op != token.AND {
		return true
	}

	// &...
	t := v.Pass.TypesInfo.TypeOf(x.X)
	if !v.zeroSizedType(t) {
		return true
	}

	message := fmt.Sprintf("address of zero-size variable of type %q", t)
	fixes := v.removeOp(x, x.X)
	v.report(x, message, fixes)

	return fixes == nil
}
