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

var _ = &B{} == &A{} // want " \\(zl:cmp\\)$" " \\(zl:add\\)$" " \\(zl:add\\)$"

func (*B) Combine(_ *B) *B { // want " \\(zl:rcv\\)$" " \\(zl:par\\)$" " \\(zl:res\\)$"
	return &B{} // want " \\(zl:add\\)$"
}

func (*B) NoLock() {} // want " \\(zl:rcv\\)$"

func Ptr[T any](v T) *T { return &v }

type greeter [5][5]struct{}

type greeterAlias = *greeter // want " \\(zl:dcl\\)$"

func (g greeterAlias) String() string { // want " \\(zl:rcv\\)$"
	return "hello, world"
}

var _ fmt.Stringer = &greeter{} // want " \\(zl:add\\)$"

var _ fmt.Stringer = (*greeter)(nil) // want " \\(zl:nil\\)$"

var _ fmt.Stringer = new(greeter) // want " \\(zl:new\\)$"

var _ = (*greeter).String(nil) // want " \\(zl:mex\\)$" " \\(zl:arg\\)$"

type greeter2[T any] [5][5][0]T

func (g *greeter2[T]) String() string { // want " \\(zl:rcv\\)$"
	return "hello, world"
}

var _ fmt.Stringer = &greeter2[int]{} // want " \\(zl:add\\)$"

type C2 struct{}

func (*C2) String() string {
	return "hello, world"
}

var _ fmt.Stringer = (*C2)(nil)

var _ = (*C2).String(nil)

type D struct{ _ int }

func (*D) String() string {
	return "hello, world"
}

var _ fmt.Stringer = (*D)(nil)

var _ = (*D).String(nil)
