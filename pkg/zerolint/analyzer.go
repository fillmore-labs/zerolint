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

package zerolint

import (
	"log"

	"fillmore-labs.com/zerolint/pkg/internal/passes/excluded"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Documentation constants.
const (
	Name = "zerolint"
	Doc  = `zerolint detects and helps fix unnecessary or problematic usage of pointers to
zero-sized types (e.g., *struct{} or *[0]byte).

Pointers to zero-size types (ZSTs) can be problematic:
- They carry very little information.
- Two pointers to distinct zero-size variables may or may not compare equal.
  This can lead to subtle bugs.
- The pointers themselves are not zero-sized and might waste memory in
  data structures, on the stack and in the CPU cache.

This analyzer helps identify such patterns to encourage using ZSTs by value or
finding alternative designs and promotes clearer, more efficient, and
spec-compliant Go code.`
)

// New creates and returns a new [analysis.Analyzer] to detect pointers to zero-length types.
func New(opts ...Option) *analysis.Analyzer {
	o := options{
		logger: log.Default(),
	}
	Options(opts).apply(&o)

	a := &analysis.Analyzer{
		Name:     Name,
		Doc:      Doc,
		URL:      "https://pkg.go.dev/fillmore-labs.com/zerolint/pkg/zerolint",
		Requires: []*analysis.Analyzer{inspect.Analyzer, excluded.Analyzer},
	}
	a.Run = o.run(&a.Flags)

	return a
}
