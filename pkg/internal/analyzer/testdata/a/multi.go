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

package a

type empty struct{}

func anonarg(*empty) { // want "function has pointer parameter to zero-sized type"
}

func twoarg(a, b *empty) { // want "function parameters \"a\", \"b\" point to zero-sized type"
	var c, d *empty = a, b // want "variables \"c\", \"d\" point to zero-sized type"

	_, _ = c, d
}

func oneresult() (e *empty) { // want "function result \"e\" points to zero-sized type"
	return &empty{} // want " \\(zl:add\\)$"
}

func tworesult() (e, f *empty) { // want "function results \"e\", \"f\" point to zero-sized type"
	return &empty{}, &empty{} // want " \\(zl:add\\)$" " \\(zl:add\\)$"
}
