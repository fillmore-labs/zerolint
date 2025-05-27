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

package set_test

import (
	"reflect"
	"slices"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/internal/set"
)

func TestSet(t *testing.T) {
	t.Parallel()

	// given
	s := New[int]()

	// when
	s.Insert(1)

	// then
	if !s.Has(1) {
		t.Error("Expected 1 to be set")
	}
}

func TestUnset(t *testing.T) {
	t.Parallel()

	// given
	s := New[int]()

	// when

	// then
	if s.Has(1) {
		t.Error("Expected 1 to be unset")
	}
}

func TestElements(t *testing.T) {
	t.Parallel()

	// given
	s := New(1)

	// when
	l := s.Elements()

	// then
	if len(l) != 1 || l[0] != 1 {
		t.Errorf("Expected l to be [1], got %v", l)
	}
}

func TestSorted(t *testing.T) {
	t.Parallel()

	type args struct {
		s Set[int]
	}

	tests := []struct {
		name string
		args args
		want []int
	}{
		{"2, 1", args{New(2, 1)}, []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Sorted(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sorted(%q) = %v, want %v", tt.args.s, got, tt.want)
			}
		})
	}
}

func TestAllSorted(t *testing.T) {
	t.Parallel()

	type args struct {
		s Set[int]
	}

	tests := []struct {
		name string
		args args
		want []int
	}{
		{"2, 1", args{New(2, 1)}, []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := slices.Collect(AllSorted(tt.args.s)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllSorted(%q) = %v, want %v", tt.args.s, got, tt.want)
			}
		})
	}
}
