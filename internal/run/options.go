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
	"fmt"
	"log"
	"regexp"

	"fillmore-labs.com/zerolint/analyzer/level"
	"fillmore-labs.com/zerolint/internal/excludes"
	"fillmore-labs.com/zerolint/internal/set"
)

// Options defines configurable parameters for the linter.
type Options struct {
	Level           level.LintLevel
	Excludes        set.Set[string]
	Generated       bool
	Regex           *regexp.Regexp
	Logger          *log.Logger
	ZeroTrace       bool
	WithFlags       bool
	ExcludeComments bool
}

// DefaultOptions returns *[Options] initialized with default values.
func DefaultOptions() *Options {
	return &Options{ // Default options
		Level:           level.Basic,
		Logger:          log.Default(),
		ExcludeComments: true,
	}
}

// ReadExcludedFile reads excluded type names from the specified file and adds them to the [Options.Excludes] set.
func (o *Options) ReadExcludedFile(name string) error {
	if name == "" {
		return nil
	}

	// If the "excluded" flag was provided, amend programmatic excludes.
	excludedTypeNames, err := excludes.ReadExcludes(osFS{}, name)
	if err != nil {
		return fmt.Errorf("error handling -excluded flag: %w", err)
	}

	if o.Excludes == nil {
		o.Excludes = set.New[string]()
	}

	for _, e := range excludedTypeNames {
		o.Excludes.Add(e)
	}

	return nil
}
