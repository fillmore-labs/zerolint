// Copyright 2024 Oliver Eikemeier. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0
package filter_test

import (
	"go/token"
	"go/types"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/internal/filter"
	"fillmore-labs.com/zerolint/pkg/internal/set"
)

func TestFilter_ExcludedType(t *testing.T) { //nolint:funlen
	t.Parallel()

	const (
		pos1 token.Pos = iota + 1
		pos2
	)

	excludedSet := set.New(pos1)
	filter := New(excludedSet)

	pkg := types.NewPackage("example.com/testpkg", "testpkg")

	typeNameExcluded := types.NewTypeName(pos1, pkg, "ExcludedType", nil)
	typeNameNotExcluded := types.NewTypeName(pos2, pkg, "NotExcludedType", nil)

	tests := []struct {
		name     string
		filter   Filter
		typeName *types.TypeName
		want     bool
	}{
		{
			name:     "Excluded type",
			filter:   filter,
			typeName: typeNameExcluded,
			want:     true,
		},
		{
			name:     "Not excluded type",
			filter:   filter,
			typeName: typeNameNotExcluded,
			want:     false,
		},
		{
			name:     "Default filter",
			filter:   Filter{},
			typeName: typeNameExcluded,
			want:     false,
		},
		{
			name:     "Nil typeName",
			filter:   filter,
			typeName: nil,
			want:     false,
		},
		{
			name:     "Default filter and nil typeName",
			filter:   Filter{},
			typeName: nil,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.filter.ExcludedType(tt.typeName); got != tt.want {
				t.Errorf("Filter.ExcludedType() = %v, want %v", got, tt.want)
			}
		})
	}
}
