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

package msg

import (
	"fmt"
	"go/ast"
	"go/types"
	"strconv"
	"strings"

	"fillmore-labs.com/zerolint/pkg/internal/diag"
)

// Formatf creates a CategorizedMessage by formatting the message
// based on the provided category, format string, and arguments.
func Formatf(cat diag.Category, valueMethod bool, format string, args ...any) diag.CategorizedMessage {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, format, args...)
	_, _ = sb.WriteString(" (zl:")
	_, _ = sb.WriteString(cat.String())

	if valueMethod {
		_ = sb.WriteByte('+')
	}

	_ = sb.WriteByte(')')

	return diag.CategorizedMessage{
		Category: cat,
		Message:  sb.String(),
	}
}

// Formatter defines an interface for generating diagnostic messages based on the context
// (struct fields, function parameters, or results) and the number of names involved.
type Formatter interface {
	// zeroMsg creates a message for unnamed items (like embedded types)
	ZeroMsg(typ types.Type, valueMethod bool) diag.CategorizedMessage

	// singularMsg creates a message for a single named item
	SingularMsg(typ types.Type, valueMethod bool, name string) diag.CategorizedMessage

	// pluralMsg creates a message for multiple named items with the same type
	PluralMsg(typ types.Type, valueMethod bool, names string) diag.CategorizedMessage
}

// FormatMessage generates a diagnostic message and category for fields, variables, or parameters.
func FormatMessage(
	formatter Formatter, typ types.Type, valueMethod bool, names []*ast.Ident,
) diag.CategorizedMessage {
	switch len(names) {
	case 0:
		return formatter.ZeroMsg(typ, valueMethod)

	case 1:
		return formatter.SingularMsg(typ, valueMethod, names[0].Name)

	default:
		quoted := make([]string, len(names))
		for i, n := range names {
			quoted[i] = strconv.Quote(n.Name)
		}

		return formatter.PluralMsg(typ, valueMethod, strings.Join(quoted, ", "))
	}
}
