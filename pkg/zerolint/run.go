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
	"flag"
	"fmt"
	"regexp"
	"sync"

	"fillmore-labs.com/zerolint/pkg/internal/excludes"
	"fillmore-labs.com/zerolint/pkg/internal/set"
	"fillmore-labs.com/zerolint/pkg/internal/visitor"
	"golang.org/x/tools/go/analysis"
)

// run returns a function that executes the analysis pass using the provided options.
// If zero-sized types are detected and zeroTrace is enabled, the function logs the detected types.
func (o options) run(s *flag.FlagSet) func(pass *analysis.Pass) (any, error) {
	delayedOptions := optionsFromFlags(&o, s)

	return func(pass *analysis.Pass) (any, error) {
		if err := delayedOptions(); err != nil {
			return nil, err
		}

		v := visitor.New(o.Options)
		res, err := v.Run(pass)

		if err == nil && o.zeroTrace && o.logger != nil && v.HasDetected() {
			o.logger.Printf("found zero-sized types in %q:\n", pass.Pkg.Path())

			for name := range v.AllDetected() {
				o.logger.Printf("- %s\n", name)
			}
		}

		return res, err
	}
}

func optionsFromFlags(o *options, s *flag.FlagSet) func() error {
	if !o.flags {
		return func() error { return nil }
	}

	var (
		regex    regexp.Regexp // For -match flag; o.Regex is *regexp.Regexp
		excluded string        // For -excluded flag; o.Excludes is []string
	)

	// Use programmatic options as defaults for flags.
	s.TextVar(&o.Level, "level", o.Level, "analysis level (Default, Extended, Full)")
	s.TextVar(&regex, "match", &regexp.Regexp{}, "only check types matching this regex, useful with -fix")
	s.StringVar(&excluded, "excluded", "", "read excluded types from this file")
	s.BoolVar(&o.zeroTrace, "zerotrace", o.zeroTrace, "trace found zero-sized types")
	s.BoolVar(&o.Generated, "generated", o.Generated, "check generated files")

	if s.Lookup("V") == nil {
		s.Var(versionFlag{}, "V", "print version and exit")
	}

	// Delayed evaluation to allow the flags to be parsed before the analysis is run.
	delayedOptions := func() error {
		// If the -excluded flag was provided, read from the file and override programmatic excludes.
		if excluded != "" {
			excludedTypeNames, err := excludes.ReadExcludes(osFS{}, excluded)
			if err != nil {
				return fmt.Errorf("error handling -excluded flag: %w", err)
			}

			o.Excludes = set.New(excludedTypeNames...)
		}

		// If -match flag was used, override programmatic regex.
		if regex.String() != "" {
			o.Regex = &regex
		}

		return nil
	}

	// Avoid re-reading the excluded file for every package
	return sync.OnceValue(delayedOptions)
}
