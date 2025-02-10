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

type GA[T any] struct {
	_ T
}

type GB[T any] = GA[T]

type Ety = struct{}

func Check(x any) {
	switch x.(type) {
	case *GB[Ety]: // want "pointer to zero-sized type"
		fmt.Println("*GB[Ety]")
	}

	switch x.(type) {
	case *GA[Ety]:
		fmt.Println("*GA[Ety]")
	}
}
