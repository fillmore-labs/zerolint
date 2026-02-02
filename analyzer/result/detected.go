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

package result

import "slices"

// Detected represents a collection of detected zero-sized types.
type Detected struct {
	detected map[string]bool
}

// New initializes and returns a [Detected] instance using the provided map of detected zero-sized types.
func New(detected map[string]bool) Detected {
	return Detected{detected: detected}
}

// Empty tells whether any zero-sized types have been detected during analysis.
func (d Detected) Empty() bool {
	return len(d.detected) == 0
}

// Sorted returns a sorted slice of all detected zero-sized types.
func (d Detected) Sorted() []string {
	sl := make([]string, len(d.detected))
	i := 0

	for n, m := range d.detected {
		if m {
			sl[i] = n + " (value methods)"
		} else {
			sl[i] = n
		}

		i++
	}

	slices.Sort(sl)

	return sl
}
