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

//nolint:gochecknoglobals
package analyzer

import (
	"log"

	"fillmore-labs.com/zerolint/pkg/analyzer/level"
	"fillmore-labs.com/zerolint/pkg/internal/excludes"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// Documentation constants.
const (
	Name = "zerolint"
	Doc  = `checks for usage of pointers to zero-length variables
	
Pointer to zero-length variables carry very little information and
can often be avoided.`
)

// Analyzer checks for usage of pointers to zero-length variables.
var Analyzer = &analysis.Analyzer{
	Name:     Name,
	Doc:      Doc,
	URL:      "https://pkg.go.dev/fillmore-labs.com/zerolint/pkg/analyzer",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var (
	// Excludes is a file containing a list of types to exclude from the analysis.
	Excludes string

	// ZeroTrace enables tracing of found zero-sized types.
	ZeroTrace bool

	// Level enables full analysis, which should be handles manually.
	Level level.LintLevel

	// Generated enables checking generated files.
	Generated bool

	// Logger to log found zero-sized types to.
	Logger *log.Logger
)

func init() { //nolint:gochecknoinits
	Analyzer.Flags.StringVar(&Excludes, "excluded", "", "read excluded types from this file")
	Analyzer.Flags.BoolVar(&ZeroTrace, "zerotrace", false, "trace found zero-sized types")
	Analyzer.Flags.TextVar(&Level, "level", level.Default, "analysis level (Default, Extended, Full)")
	Analyzer.Flags.BoolVar(&Generated, "generated", false, "check generated files")
}

// run applies the analyzer to a package.
func run(pass *analysis.Pass) (any, error) {
	// Read the list of excluded types from the file specified by the "excluded" flag.
	ex, err := excludes.ReadExcludes(osFS{}, Excludes)
	if err != nil {
		return nil, err
	}

	return NewRun(
		WithExcludes(ex),
		WithZeroTrace(ZeroTrace),
		WithLevel(Level),
		WithGenerated(Generated),
		WithLogger(Logger),
	)(pass)
}
