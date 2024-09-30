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

	"golang.org/x/tools/go/analysis"
)

// visitFunc checks if the function declaration has a receiver of a pointer to a zero-sized type.
func (v Visitor) visitFunc(x *ast.FuncDecl) bool {
	// Only process methods.
	if x.Recv == nil || len(x.Recv.List) != 1 {
		return true
	}

	recv := x.Recv.List[0]
	recvType := v.TypesInfo.TypeOf(recv.Type)
	elem, ok := v.zeroSizedTypePointer(recvType)
	if !ok { // Not a pointer receiver or no pointer to a zero-sized type.
		return true
	}

	var fixes []analysis.SuggestedFix
	if s, ok := recv.Type.(*ast.StarExpr); ok {
		fixes = v.removeOp(s, s.X)
	}

	message := fmt.Sprintf("method receiver is pointer to zero-size variable of type %q", elem)
	v.report(recv, message, fixes)

	return true
}
