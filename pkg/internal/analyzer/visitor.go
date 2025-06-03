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

	"fillmore-labs.com/zerolint/pkg/internal/checker"
	"fillmore-labs.com/zerolint/pkg/internal/diag"
	"fillmore-labs.com/zerolint/pkg/internal/set"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

// visitor is an AST visitor for analyzing the usage of pointers to zero-sized types.
// It identifies various patterns where such pointers might be used unnecessarily.
type visitor struct {
	check     checker.Checker // Helper for analyzing.
	diag      diag.Diag       // Helper for reporting.
	level     level.LintLevel // Analysis level.
	generated bool            // Analyze generated source, too.

	// Tracks *[ast.StarExpr] positions that have already been processed to avoid duplicate diagnostics or fixes.
	seenStars set.Set[token.Pos]
}
