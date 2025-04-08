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

package checker

import (
	"fmt"
	"go/ast"
	"go/types"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// MsgFormatter defines an interface for generating diagnostic messages based on the context
// (struct fields, function parameters, or results) and the number of names involved.
type MsgFormatter[C fmt.Stringer] interface {
	// ZeroMsg creates a message for unnamed fields (like embedded types)
	ZeroMsg(t types.Type) (C, string)

	// SingularMsg creates a message for a single named field
	SingularMsg(name string, t types.Type) (C, string)

	// PluralMsg creates a message for multiple named fields with the same type
	PluralMsg(names string, t types.Type) (C, string)
}

// FormatMessage generates a diagnostic message and category for fields, variables, or parameters.
func FormatMessage[C fmt.Stringer](names []*ast.Ident, msgFormatter MsgFormatter[C], elem types.Type) (C, string) {
	switch len(names) {
	case 0:
		return msgFormatter.ZeroMsg(elem)

	case 1:
		return msgFormatter.SingularMsg(names[0].Name, elem)

	default:
		quoted := make([]string, len(names))
		for i, n := range names {
			quoted[i] = strconv.Quote(n.Name)
		}

		return msgFormatter.PluralMsg(strings.Join(quoted, ", "), elem)
	}
}

// Report adds a diagnostic message to the analysis pass results.
func (c *Checker) Report(
	rng analysis.Range, cat fmt.Stringer, plus bool, message string, fixes []analysis.SuggestedFix,
) {
	var p string
	if plus {
		p = "+"
	}

	catString := cat.String()
	msg := fmt.Sprintf("%s (zl:%s%s)", message, catString, p)
	c.pass.Report(analysis.Diagnostic{
		Pos:            rng.Pos(),
		End:            rng.End(),
		Category:       catString,
		Message:        msg,
		URL:            "", // "https://blog.fillmore-labs.com/posts/zerolint" + "#" + catString,
		SuggestedFixes: fixes,
	})
}
