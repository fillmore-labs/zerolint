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

package analyzer

import (
	"fillmore-labs.com/zerolint/pkg/visitor"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const (
	Name = "zerolint"
	Doc  = `checks for usage of pointers to zero-length variables
	
Pointer to zero-length variables carry very little information and
can be avoided in most cases.`
)

var Analyzer = &analysis.Analyzer{ //nolint:gochecknoglobals
	Name:     Name,
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func init() { //nolint:gochecknoinits
	Analyzer.Flags.StringVar(&Excludes, "excluded", "", "read excluded types from this file")
	Analyzer.Flags.BoolVar(&ZeroTrace, "zerotrace", false, "trace found zero-sized types")
}

var ZeroTrace bool //nolint:gochecknoglobals

func run(pass *analysis.Pass) (any, error) {
	excludes, err := ReadExcludes()
	if err != nil {
		return nil, err
	}

	v := visitor.Visitor{Pass: pass, Excludes: excludes, ZeroTrace: ZeroTrace}
	v.Run()

	return any(nil), nil
}
