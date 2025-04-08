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
	"log"
	"regexp"

	"fillmore-labs.com/zerolint/pkg/internal/excludes"
	"fillmore-labs.com/zerolint/pkg/internal/set"
	"fillmore-labs.com/zerolint/pkg/internal/visitor"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
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

		if err == nil && o.zeroTrace && v.HasDetected() {
			logger := o.logger
			if logger == nil {
				logger = log.Default()
			}

			logger.Printf("found zero-sized types in %q:\n", pass.Pkg.Path())

			for name := range v.AllDetected() {
				logger.Printf("- %s\n", name)
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
		lvl   level.LintLevel
		match regexp.Regexp
	)

	s.TextVar(&lvl, "level", level.Default, "analysis level (Default, Extended, Full)")
	s.TextVar(&match, "match", &regexp.Regexp{}, "only check types matching this regex, useful with -fix")
	excluded := s.String("excluded", "", "read excluded types from this file")
	zeroTrace := s.Bool("zerotrace", false, "trace found zero-sized types")
	generated := s.Bool("generated", false, "check generated files")

	if s.Lookup("V") == nil {
		s.Var(versionFlag{}, "V", "print version and exit")
	}

	// Delayed evaluation to allow the flags to be parsed before the analysis is run.
	return func() error {
		// Read the list of excluded types from the file specified by the "excluded" flag.
		ex, err := excludes.ReadExcludes(osFS{}, *excluded)
		if err != nil {
			return err
		}

		o.Level = lvl
		o.Excludes = set.New(ex...)
		o.zeroTrace = *zeroTrace
		o.Generated = *generated

		if match.String() != "" {
			o.Regex = &match
		}

		return nil
	}
}
