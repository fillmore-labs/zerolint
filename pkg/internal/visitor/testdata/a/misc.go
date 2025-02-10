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

import "fmt"

func Undiagnosed() {
	i := 1 + 1
	i = *&i

	f := func(int) {}
	f(i)

	_ = !false
}

type A [0]string

type B = A

var _ = &B{} == &A{} // want "comparison of pointers to zero-size variables" "address of zero-size variable" "address of zero-size variable"

func (*B) Combine(_ *B) *B { // want "pointer to zero-sized type" "pointer to zero-sized type" "pointer to zero-sized type"
	return &B{} // want "address of zero-size variable"
}

func Ptr[T any](v T) *T { return &v }

type greeter [5][5]struct{}

type greeterAlias = *greeter // want "pointer to zero-sized type"

func (g greeterAlias) String() string {
	return "hello, world"
}

type generatedAlias = *generated

var _ fmt.Stringer = &greeter{} // want "address of zero-size variable"

var _ fmt.Stringer = (*greeter)(nil) // want "cast of nil to pointer to zero-size variable"

var _ fmt.Stringer = new(greeter) // want "new called on zero-sized type"

type greeter2[T any] [5][5][0]T

func (g *greeter2[T]) String() string { // want "pointer to zero-sized type"
	return "hello, world"
}

var _ fmt.Stringer = &greeter2[int]{} // want "address of zero-size variable"

type C2 struct{}

func (*C2) String() string {
	return "hello, world"
}

var _ fmt.Stringer = (*C2)(nil)

type D struct{ _ int }

func (*D) String() string {
	return "hello, world"
}

var _ fmt.Stringer = (*D)(nil)
