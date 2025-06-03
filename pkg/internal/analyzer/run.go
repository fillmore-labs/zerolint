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
	"errors"
	"go/token"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"fillmore-labs.com/zerolint/pkg/internal/checker"
	"fillmore-labs.com/zerolint/pkg/internal/filter"
	"fillmore-labs.com/zerolint/pkg/internal/passes/exclusions"
	"fillmore-labs.com/zerolint/pkg/internal/set"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
	"fillmore-labs.com/zerolint/pkg/zerolint/result"
)

// ErrNoInspectorResult is returned when the ast inspector is missing.
var ErrNoInspectorResult = errors.New("zerolint: inspector result missing")

// Options defines configurable parameters for the analyzer.
type Options struct {
	Level     level.LintLevel
	Excludes  set.Set[string]
	Generated bool
	Regex     *regexp.Regexp
}

// Run performs the actual analysis on the provided [analysis.Pass].
func (o Options) Run(pass *analysis.Pass) (any, error) {
	v := &visitor{
		check: checker.Checker{
			Excludes: o.Excludes,
		},
		level:     o.Level,
		generated: o.Generated,
	}
	if o.Regex != nil && o.Regex.String() != "" {
		v.check.Regex = o.Regex
	}

	v.check.Prepare()
	v.diag.Prepare(pass)
	v.seenStars = make(set.Set[token.Pos])

	if excludedTypeDefs, err := exclusions.CalculateExclusions(pass); err == nil {
		v.check.ExcludedTypeDefs = filter.New(excludedTypeDefs)
	} else if !errors.Is(err, exclusions.ErrNoExclusionsResult) {
		return nil, err
	}

	in, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, ErrNoInspectorResult
	}

	types := v.nodeFilter()
	in.Root().Inspect(types, v.dispatch)

	return result.New(v.check.Detected), nil
}
