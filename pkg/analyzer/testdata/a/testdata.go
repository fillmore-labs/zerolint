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

import (
	"errors"
	"fmt"

	xerrors "golang.org/x/exp/errors"
)

type empty struct{}

type typedError[T any] struct {
	_ [0]T
}

func (*typedError[_]) Error() string { // want "pointer to zero-size type"
	return "an error"
}

var (
	_      error = &typedError[any]{}         // want "address of zero-size variable"
	ErrOne       = &(typedError[int]{})       // want "address of zero-size variable"
	ErrTwo       = (new)(typedError[float64]) // want "new called on zero-size type"
)

type myErrors struct{}

func (myErrors) Is(err, target error) bool {
	return false
}

var myErrs = myErrors{}

func Exported() {
	var x [0]string
	var y [0]string

	if errors.Is(ErrOne, nil) {
		fmt.Println("nil")
	}

	func() {
		errors := myErrs
		if errors.Is(ErrOne, ErrTwo) {
			fmt.Println("one is two")
		}
	}()

	if xerrors.Is(func() error { // want "comparison of pointer to zero-size variable"
		return ErrOne
	}(), ErrTwo) {
		fmt.Println("equal")
	}

	var err *typedError[int] // want "pointer to zero-size type"
	_ = errors.As(ErrOne, &err)

	_ = (new)(struct{}) // want "new called on zero-size type"

	_ = new(empty) // want "new called on zero-size type"

	xp, yp := &x, &y // want "address of zero-size variable" "address of zero-size variable"

	_ = *xp // want "pointer to zero-size variable"

	if xp == yp { // want "comparison of pointers to zero-size variables"
		fmt.Println("equal")
	}

	if xp != yp { // want "comparison of pointers to zero-size variables"
		fmt.Println("not equal")
	}

	if xp == nil {
		fmt.Println("nil")
	}

	_, _ = any(xp).((*[0]string)) // want "pointer to zero-size type"

	switch any(xp).(type) {
	case (*[0]string): // want "pointer to zero-size type"
	case string:
	}
}

func Undiagnosed() {
	i := 1 + 1
	i = *&i

	f := func(int) {}
	f(i)

	_ = !false
}

type A [0]string

type B = A

func (*B) Combine(_ *B) *B { // want "pointer to zero-size type" "pointer to zero-size type" "pointer to zero-size type"
	return &B{} // want "address of zero-size variable"
}

func Ptr[T any](v T) *T { return &v }

type greeter [5][5]struct{}

type greeterAlias = *greeter // want "pointer to zero-size type"

func (g greeterAlias) String() string {
	return "hello, world"
}

var _ fmt.Stringer = &greeter{} // want "address of zero-size variable"

var _ fmt.Stringer = (*greeter)(nil) // want "cast of nil to pointer to zero-size variable"

var _ fmt.Stringer = new(greeter) // want "new called on zero-size type"

type greeter2[T any] [5][5][0]T

func (g *greeter2[T]) String() string { // want "pointer to zero-size type"
	return "hello, world"
}

var _ fmt.Stringer = &greeter2[int]{} // want "address of zero-size variable"

type C struct{}

func (*C) String() string {
	return "hello, world"
}

var _ fmt.Stringer = (*C)(nil)

type D struct{ _ int }

func (*D) String() string {
	return "hello, world"
}

var _ fmt.Stringer = (*D)(nil)
