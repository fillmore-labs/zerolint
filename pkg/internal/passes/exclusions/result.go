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

package exclusions

import (
	"errors"
	"go/token"

	"golang.org/x/tools/go/analysis"

	"fillmore-labs.com/zerolint/pkg/internal/set"
)

type exclusionsResult struct {
	facts []analysis.ObjectFact
}

func (p pass) newResult() exclusionsResult {
	return exclusionsResult{facts: p.AllObjectFacts()}
}

// ErrNoExclusionsResult is returned when the [Analyzer]s result is missing from the [analysis.Pass].
var ErrNoExclusionsResult = errors.New("result of exclusions.Analyzer missing")

// ResultOf retrieves a set of token positions for type definitions that have been excluded by the [Analyzer].
// It returns the set of excluded positions or an error if the exclusion results are not available.
func ResultOf(pass *analysis.Pass) (set.Set[token.Pos], error) {
	excludedResult, ok := pass.ResultOf[Analyzer].(exclusionsResult)
	if !ok {
		return nil, ErrNoExclusionsResult
	}

	excludedTypeDefs := set.New[token.Pos]()
	for obj := range AllFacts[*excludedFact](excludedResult.facts) {
		excludedTypeDefs.Add(obj.Pos())
	}

	return excludedTypeDefs, nil
}
