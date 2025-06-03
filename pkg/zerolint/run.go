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
	"flag"
	"fmt"
	"regexp"

	"golang.org/x/tools/go/analysis"

	"fillmore-labs.com/zerolint/pkg/internal/excludes"
	"fillmore-labs.com/zerolint/pkg/internal/set"
	"fillmore-labs.com/zerolint/pkg/zerolint/result"
)

// run is the function that executes an analysis pass using the provided options.
// If zero-sized types are detected and zeroTrace is enabled, the function logs the detected types.
func (o *options) run(pass *analysis.Pass) (any, error) {
	// Avoid re-reading the excluded file for every package
	if o.excludeRead.Do(o.readExcludedFile); o.excludeReadErr != nil {
		return nil, o.excludeReadErr
	}

	res, err := o.analyzer.Run(pass)

	d, ok := res.(result.Detected)
	if ok && err == nil && o.zeroTrace && o.logger != nil && !d.Empty() {
		o.logger.Printf("Found zero-sized types in %q:\n", pass.Pkg.Path())

		for _, name := range d.Sorted() {
			o.logger.Printf("- %s\n", name)
		}
	}

	return res, err
}

func (o *options) readExcludedFile() {
	if o.excludedFile == "" {
		return
	}

	// If the -excluded flag was provided, amend programmatic excludes.
	excludedTypeNames, err := excludes.ReadExcludes(osFS{}, o.excludedFile)
	if err != nil {
		o.excludeReadErr = fmt.Errorf("error handling -excluded flag: %w", err)

		return
	}

	if o.analyzer.Excludes == nil {
		o.analyzer.Excludes = set.New[string]()
	}

	for _, e := range excludedTypeNames {
		o.analyzer.Excludes.Insert(e)
	}
}

// flags returns a [flag.FlagSet] containing command-line flags that can
// configure the analyzer's behavior. These flags correspond to the fields
// in the [options] struct.
// The analysis driver uses the returned FlagSet to pass command-line
// arguments to the analyzer.
func (o *options) flags() flag.FlagSet {
	var flags flag.FlagSet
	if !o.withFlags {
		return flags
	}

	if o.analyzer.Regex == nil {
		o.analyzer.Regex = &regexp.Regexp{}
	}

	// Use programmatic options as defaults for flags.
	flags.TextVar(&o.analyzer.Level, "level", o.analyzer.Level, "analysis level (Default, Extended, Full)")
	flags.TextVar(o.analyzer.Regex, "match", o.analyzer.Regex, "only check types matching this regex, useful with -fix")
	flags.StringVar(&o.excludedFile, "excluded", o.excludedFile, "read excluded types from this file")
	flags.BoolVar(&o.zeroTrace, "zerotrace", o.zeroTrace, "trace found zero-sized types")
	flags.BoolVar(&o.analyzer.Generated, "generated", o.analyzer.Generated, "check generated files")

	return flags
}
