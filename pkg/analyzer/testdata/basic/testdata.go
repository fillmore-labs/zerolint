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

package basic

import (
	"errors"
	"fmt"
)

type myError struct{}

func (*myError) Error() string { // want "method receiver is pointer to zero-size variable"
	return "my error"
}

var (
	ErrOne = &myError{}
	ErrTwo = new(myError)
)

func Exported() {
	if errors.Is(nil, ErrOne) {
		fmt.Println("nil")
	}

	if errors.Is(func() error { // want "comparison of pointer to zero-size variable"
		return ErrOne
	}(), ErrTwo) {
		fmt.Println("equal")
	}

	var err *myError
	_ = errors.As(ErrOne, &err)

	if ErrOne == ErrTwo { // want "comparison of pointers to zero-size variables"
		fmt.Println("equal")
	}

	if ErrOne != ErrTwo { // want "comparison of pointers to zero-size variables"
		fmt.Println("not equal")
	}
}

type D struct{ _ int }

func (*D) String() string {
	return "hello, world"
}

var _ fmt.Stringer = (*D)(nil)
