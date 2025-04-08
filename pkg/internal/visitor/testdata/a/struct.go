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

package a

type typedError[T any] struct {
	_ [0]T
}

type embeddedPointer struct {
	*empt       // want "embedded pointer to zero-sized type"
	t     *empt // want "field t points to zero-sized type"
	u, v  *empt // want "fields u, v point to zero-sized type"
}

func (*typedError[_]) Error() string { // want "error interface implemented on pointer to zero-sized type"
	return "an error"
}

var (
	_      error = &typedError[any]{}         // want "address of zero-size variable"
	ErrOne       = &(typedError[int]{})       // want "address of zero-size variable"
	ErrTwo       = (new)(typedError[float64]) // want "new called on zero-sized type"
)
