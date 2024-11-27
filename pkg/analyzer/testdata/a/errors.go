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

type myErrors struct{}

func (myErrors) Is(_, _ error) bool {
	return false
}

var myErrs = myErrors{}

func IgnoreErrors() {
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

	var err *typedError[int] // want "pointer to zero-sized type"
	_ = errors.As(ErrOne, &err)
}
