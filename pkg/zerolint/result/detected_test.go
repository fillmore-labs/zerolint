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

package result_test

import (
	"slices"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/zerolint/result"
)

func TestResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   map[string]bool
		wantLen int
	}{
		{
			name:    "empty map",
			input:   map[string]bool{},
			wantLen: 0,
		},
		{
			name: "non-empty map",
			input: map[string]bool{
				"TypeA": false,
				"TypeB": true,
			},
			wantLen: 2,
		},
		{
			name:    "nil map",
			input:   nil,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := New(tt.input)

			if got, want := d.Empty(), tt.wantLen == 0; got != want {
				t.Errorf("Expected Empty() to be %t, got %t", want, got)
			}

			sorted := d.Sorted()

			if got, want := len(sorted), tt.wantLen; got != want {
				t.Errorf("Expected Sorted() to return slice of length %d, got %d", want, got)
			}

			if !slices.IsSorted(sorted) {
				t.Errorf("Sorted() returned slice is not sorted: %v", sorted)
			}
		})
	}
}
