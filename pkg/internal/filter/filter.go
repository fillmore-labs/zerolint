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

package filter

import (
	"go/token"
	"go/types"

	"fillmore-labs.com/zerolint/pkg/internal/set"
)

// Filter holds a set of token positions for type definitions that have been
// marked as excluded.
type Filter struct {
	excludedTypeDefs set.Set[token.Pos]
}

// New creates a new Filter with the specified set of excluded type definitions.
func New(excludedTypeDefs set.Set[token.Pos]) Filter {
	return Filter{excludedTypeDefs: excludedTypeDefs}
}

// ExcludedType checks if a given [types.TypeName] (representing a defined type)
// has been marked as excluded.
func (r Filter) ExcludedType(tn *types.TypeName) bool {
	if tn == nil {
		return false
	}

	return r.excludedTypeDefs.Has(tn.Pos())
}
