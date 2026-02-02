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

package run

import (
	"errors"

	"golang.org/x/tools/go/analysis"

	"fillmore-labs.com/zerolint/analyzer/result"
	"fillmore-labs.com/zerolint/internal/analyzer"
	"fillmore-labs.com/zerolint/internal/checker"
)

// ErrNoInspectorResult is returned when the ast inspector is missing.
var ErrNoInspectorResult = errors.New("zerolint: inspector result missing")

// Run is the function that executes an analysis pass using the provided options.
// If zero-sized types are detected and zeroTrace is enabled, the function logs the detected types.
func (o *Options) Run(pass *analysis.Pass) (any, error) {
	v := &analyzer.Visitor{
		Check: checker.Checker{
			Excludes: o.Excludes,
		},
		Level:     o.Level,
		Generated: o.Generated,
	}
	if o.Regex != nil && o.Regex.String() != "" {
		v.Check.Regex = o.Regex
	}

	res, err := v.Run(pass)
	if err != nil {
		return nil, err
	}

	d, ok := res.(result.Detected)
	if ok && o.ZeroTrace && o.Logger != nil && !d.Empty() {
		o.Logger.Printf("Found zero-sized types in %q:\n", pass.Pkg.Path())

		for _, name := range d.Sorted() {
			o.Logger.Printf("- %s\n", name)
		}
	}

	return d, nil
}
