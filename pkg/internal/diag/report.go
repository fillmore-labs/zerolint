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

package diag

import "golang.org/x/tools/go/analysis"

// Category represents an internal category code used to categorize
// different types of issues found by the linter.
type Category string

func (c Category) String() string {
	return string(c)
}

// CategorizedMessage is a message that has a category and a pre-formatted message string.
type CategorizedMessage struct {
	Message  string
	Category Category
}

// Report adds a diagnostic message to the analysis pass results using the [analysis.Pass]'s Report method.
func (d *Diag) Report(rng analysis.Range, msg CategorizedMessage, fixes []analysis.SuggestedFix) {
	d.pass.Report(analysis.Diagnostic{
		Pos:            rng.Pos(),
		End:            rng.End(),
		Category:       msg.Category.String(),
		Message:        msg.Message,
		SuggestedFixes: fixes,
		// URL:            "https://blog.fillmore-labs.com/posts/zerolint" + "#" + msg.Category,
	})
}
