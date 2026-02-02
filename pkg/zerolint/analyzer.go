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

package zerolint

import (
	"fmt"
	"reflect"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"fillmore-labs.com/zerolint/pkg/internal/excludes"
	"fillmore-labs.com/zerolint/pkg/internal/passes/exclusions"
	"fillmore-labs.com/zerolint/pkg/internal/set"
	"fillmore-labs.com/zerolint/pkg/zerolint/result"
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
	URL = "https://pkg.go.dev/fillmore-labs.com/zerolint/pkg/zerolint"
)

// New creates and returns a new [analysis.Analyzer] to detect pointers to zero-length types.
func New(opts ...Option) *analysis.Analyzer {
	o := makeOptions(opts)

	requires := []*analysis.Analyzer{inspect.Analyzer, exclusions.Analyzer}
	if !o.excludeComments {
		requires = []*analysis.Analyzer{inspect.Analyzer}
	}

	a := &analysis.Analyzer{
		Name: Name,
		Doc:  Doc,
		URL:  URL,
		Run:  o.run,

		Requires:   requires,
		ResultType: reflect.TypeFor[result.Detected](),
	}

	if o.withFlags {
		if o.regex == nil {
			o.regex = &regexp.Regexp{}
		}

		// Use programmatic options as defaults for flags.
		a.Flags.TextVar(&o.level, "level", o.level, "analysis `level` (basic, extended, full)")
		a.Flags.TextVar(o.regex, "match", o.regex, "only check types matching this `regex`, useful with -fix")
		a.Flags.Func("excluded", "read excluded types from this `file`", o.readExcludedFile)
		a.Flags.BoolVar(&o.zeroTrace, "zerotrace", o.zeroTrace, "trace found zero-sized types")
		a.Flags.BoolVar(&o.generated, "generated", o.generated, "check generated files")
	}

	return a
}

func (o *options) readExcludedFile(name string) error {
	if name == "" {
		return nil
	}

	// If the -excluded flag was provided, amend programmatic excludes.
	excludedTypeNames, err := excludes.ReadExcludes(osFS{}, name)
	if err != nil {
		return fmt.Errorf("error handling -excluded flag: %w", err)
	}

	if o.excludes == nil {
		o.excludes = set.New[string]()
	}

	for _, e := range excludedTypeNames {
		o.excludes.Add(e)
	}

	return nil
}
