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

type m1 struct{}

func (m *m1) f() *m1 { return m } // want "pointer to zero-sized type" "pointer to zero-sized type"

var _ = (*m1).f(&m1{}) // want "method expression receiver is pointer to zero-size variable" "pointer to zero-sized type" "address of zero-size variable"

type m2 struct{} //zerolint:exclude

func (m *m2) f() *m2 { return m }

var _ = (*m2).f(nil) // Would be broken by '-fix'.
