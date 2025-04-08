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
	"iter"
	"regexp"

	"fillmore-labs.com/zerolint/pkg/internal/checker"
	"fillmore-labs.com/zerolint/pkg/internal/set"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

// Options defines configurable parameters for the linter.
type Options struct {
	Level     level.LintLevel
	Excludes  set.Set[string]
	Generated bool
	Regex     *regexp.Regexp
}

// Visitor is an AST visitor for analyzing the usage of pointers to zero-sized types.
// It identifies various patterns where such pointers might be used unnecessarily.
type Visitor struct {
	check     checker.Checker
	level     level.LintLevel
	generated bool
}

// New creates a new [Visitor] configured with the provided [Options].
func New(opt Options) *Visitor {
	return &Visitor{
		check: checker.Checker{
			Excludes: opt.Excludes,
			Regex:    opt.Regex,
		},
		level:     opt.Level,
		generated: opt.Generated,
	}
}

// HasDetected tells whether any zero-sized types have been detected during analysis.
func (v *Visitor) HasDetected() bool {
	return len(v.check.Detected) > 0
}

// AllDetected returns a sorted iterator over all detected zero-sized types.
func (v *Visitor) AllDetected() iter.Seq[string] {
	return set.AllSorted(v.check.Detected)
}
