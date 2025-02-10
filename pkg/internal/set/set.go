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

package set

import (
	"cmp"
	"iter"
	"slices"
)

// Set is a collection of unique elements of type T.
// It is implemented using a map with values of an empty struct.
type Set[T comparable] map[T]struct{}

// New returns a new set of type T containing elems.
func New[T comparable](elems ...T) Set[T] {
	s := make(Set[T], len(elems))
	for _, e := range elems {
		s.Insert(e)
	}

	return s
}

// Insert adds the element t to the set s.
// If t is already in s, it does nothing.
func (s Set[T]) Insert(t T) {
	s[t] = struct{}{}
}

// Has returns true if the element t is in the set s, false otherwise.
func (s Set[T]) Has(t T) bool {
	_, ok := s[t]

	return ok
}

// Elements returns the elements of s as a list.
func (s Set[T]) Elements() []T {
	sl := make([]T, len(s))
	i := 0
	for n := range s {
		sl[i] = n
		i++
	}

	return sl
}

// Sorted returns the elements of s as a sorted list.
func Sorted[T cmp.Ordered](s Set[T]) []T {
	ret := s.Elements()
	slices.Sort(ret)

	return ret
}

// AllSorted returns an iterator over the elements of s in sorted order.
func AllSorted[T cmp.Ordered](s Set[T]) iter.Seq[T] {
	return slices.Values(Sorted(s))
}
