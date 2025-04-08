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

import "golang.org/x/tools/go/analysis"

// CategorizedMessage is a message that has a category and a pre-formatted message string.
type CategorizedMessage struct {
	Message  string
	Category string
}

// Report adds a diagnostic message to the analysis pass results using the [analysis.Pass]'s Report method.
func (c *Checker) Report(rng analysis.Range, msg CategorizedMessage, fixes []analysis.SuggestedFix) {
	c.pass.Report(analysis.Diagnostic{
		Pos:            rng.Pos(),
		End:            rng.End(),
		Category:       msg.Category,
		Message:        msg.Message,
		URL:            "", // "https://blog.fillmore-labs.com/posts/zerolint" + "#" + msg.Category,
		SuggestedFixes: fixes,
	})
}
