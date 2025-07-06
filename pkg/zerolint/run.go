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
	"errors"

	"golang.org/x/tools/go/analysis"

	"fillmore-labs.com/zerolint/pkg/internal/analyzer"
	"fillmore-labs.com/zerolint/pkg/internal/checker"
	"fillmore-labs.com/zerolint/pkg/zerolint/result"
)

// ErrNoInspectorResult is returned when the ast inspector is missing.
var ErrNoInspectorResult = errors.New("zerolint: inspector result missing")

// run is the function that executes an analysis pass using the provided options.
// If zero-sized types are detected and zeroTrace is enabled, the function logs the detected types.
func (o *options) run(pass *analysis.Pass) (any, error) {
	v := &analyzer.Visitor{
		Check: checker.Checker{
			Excludes: o.excludes,
		},
		Level:     o.level,
		Generated: o.generated,
	}
	if o.regex != nil && o.regex.String() != "" {
		v.Check.Regex = o.regex
	}

	res, err := v.Run(pass)
	if err != nil {
		return nil, err
	}

	d, ok := res.(result.Detected)
	if ok && o.zeroTrace && o.logger != nil && !d.Empty() {
		o.logger.Printf("Found zero-sized types in %q:\n", pass.Pkg.Path())

		for _, name := range d.Sorted() {
			o.logger.Printf("- %s\n", name)
		}
	}

	return d, nil
}
