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

package checker

import (
	"regexp"

	"golang.org/x/tools/go/types/typeutil"

	"fillmore-labs.com/zerolint/pkg/internal/filter"
	"fillmore-labs.com/zerolint/pkg/internal/set"
)

// Checker provides helper functions for analyzing pointers to zero-sized types.
type Checker struct {
	// Cached results of zero-sized checks, used by [Checker.ZeroSizedType] to optimize repeated lookups.
	cache typeutil.Map

	// Type definitions excluded via `//nolint:zerolint` directives.
	ExcludedTypeDefs filter.Filter

	Excludes set.Set[string]

	Detected map[string]bool

	// Filter for zero-sized checks, used in [Checker.ZeroSizedType].
	Regex *regexp.Regexp
}

// Prepare initializes the [Checker] with the provided [analysis.Pass], preparing for new analysis.
func (c *Checker) Prepare() {
	if c.Excludes == nil {
		c.Excludes = set.New[string]()
	}

	c.Detected = make(map[string]bool)
}
