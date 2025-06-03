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

type m1 struct{}

func (m *m1) f(n *m1) (*m1, *m1) { return m, n } // want " \\(zl:rcv\\)$" " \\(zl:par\\)$" " \\(zl:res\\)$" " \\(zl:res\\)$"

var _, _ = (*m1).f(&m1{}, &m1{}) // want " \\(zl:mex\\)$" " \\(zl:add\\)$" " \\(zl:add\\)$"

var _, _ = (*m1).f(nil, nil) // want " \\(zl:mex\\)$" " \\(zl:arg\\)$"  " \\(zl:arg\\)$"

type m2 = m1

func (m *m2) f2(n *m2) (*m2, *m2) { return m, n } // want " \\(zl:rcv\\)$" " \\(zl:par\\)$" " \\(zl:res\\)$" " \\(zl:res\\)$"

var _, _ = (*m2).f2(nil, nil) // want " \\(zl:mex\\)$" " \\(zl:arg\\)$"  " \\(zl:arg\\)$"

type m3 struct{}

func (m m3) f(n m3) (m3, m3) { return m, n }

var _, _ = (m3).f(m3{}, m3{})

type m4 = m3

func (m m4) f2(n m4) (m4, m4) { return m, n }

var _, _ = (m4).f2(m4{}, m4{})

type m5 struct{ m4 }

func (m m5) f3(n m5) (m5, m5) { return m, n }

var _, _ = (*m5).f3(&m5{}, m5{}) // want " \\(zl:mex\\+\\)$" " \\(zl:add\\+\\)$"

var _, _ = (m5).f(m5{}, m3{})

var _, _ = (*m5).f2(&m5{}, m4{}) // want " \\(zl:mex\\+\\)$" " \\(zl:add\\+\\)$"
