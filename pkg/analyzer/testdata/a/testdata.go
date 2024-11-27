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
	"encoding/json"
	"errors"
	"fmt"

	xerrors "golang.org/x/exp/errors"
)

type empty struct{}

type typedError[T any] struct {
	_ [0]T
}

func (*typedError[_]) Error() string { // want "pointer to zero-sized type"
	return "an error"
}

var (
	_      error = &typedError[any]{}         // want "address of zero-size variable"
	ErrOne       = &(typedError[int]{})       // want "address of zero-size variable"
	ErrTwo       = (new)(typedError[float64]) // want "new called on zero-sized type"
)

type myErrors struct{}

func (myErrors) Is(_, _ error) bool {
	return false
}

var myErrs = myErrors{}

func Exported() {
	var x [0]string
	var y [0]string

	if errors.Is(ErrOne, nil) {
		fmt.Println("nil")
	}

	var oneErr *typedError[int] // want "pointer to zero-sized type"
	if errors.As(ErrOne, &oneErr) {
		fmt.Println("ErrOne is typedError[int]")
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

	_ = json.Unmarshal(nil, &myErrs)
	_ = (*json.Decoder)(nil).Decode(&myErrs)
	_ = json.NewDecoder(nil).Decode(&myErrs)

	var err *typedError[int] // want "pointer to zero-sized type"
	_ = errors.As(ErrOne, &err)

	_ = (new)(struct{}) // want "new called on zero-sized type"

	_ = new(empty) // want "new called on zero-sized type"

	xp, yp := &x, &y // want "address of zero-size variable" "address of zero-size variable"

	_ = *xp // want "pointer to zero-size variable"

	if xp == yp { // want "comparison of pointers to zero-size variables"
		fmt.Println("equal")
	}

	if xp != yp { // want "comparison of pointers to zero-size variables"
		fmt.Println("not equal")
	}

	_, _ = any(xp).((*[0]string)) // want "pointer to zero-sized type"

	switch any(xp).(type) {
	case (*[0]string): // want "pointer to zero-sized type"
	case string:
	}
}
