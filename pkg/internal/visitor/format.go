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
	"strconv"
	"strings"

	"fillmore-labs.com/zerolint/pkg/internal/checker"
)

func msgFormatf(cat category, valueMethod bool, format string, args ...any) checker.CategorizedMessage {
	var p string
	if valueMethod {
		p = "+"
	}

	const (
		extraFormat = " (zl:%s%s)"
		extraArgs   = 2
	)

	newFormat := format + extraFormat
	newArgs := make([]any, 0, len(args)+extraArgs)
	newArgs = append(newArgs, args...)
	newArgs = append(newArgs, cat, p)

	return checker.CategorizedMessage{
		Category: cat.String(),
		Message:  fmt.Sprintf(newFormat, newArgs...),
	}
}

// msgFormatter defines an interface for generating diagnostic messages based on the context
// (struct fields, function parameters, or results) and the number of names involved.
type msgFormatter interface {
	// zeroMsg creates a message for unnamed items (like embedded types)
	zeroMsg(typ types.Type, valueMethod bool) checker.CategorizedMessage

	// singularMsg creates a message for a single named item
	singularMsg(typ types.Type, valueMethod bool, name string) checker.CategorizedMessage

	// pluralMsg creates a message for multiple named items with the same type
	pluralMsg(typ types.Type, valueMethod bool, names string) checker.CategorizedMessage
}

// formatMessage generates a diagnostic message and category for fields, variables, or parameters.
func formatMessage(
	formatter msgFormatter, typ types.Type, valueMethod bool, names []*ast.Ident,
) checker.CategorizedMessage {
	switch len(names) {
	case 0:
		return formatter.zeroMsg(typ, valueMethod)

	case 1:
		return formatter.singularMsg(typ, valueMethod, names[0].Name)

	default:
		quoted := make([]string, len(names))
		for i, n := range names {
			quoted[i] = strconv.Quote(n.Name)
		}

		return formatter.pluralMsg(typ, valueMethod, strings.Join(quoted, ", "))
	}
}
