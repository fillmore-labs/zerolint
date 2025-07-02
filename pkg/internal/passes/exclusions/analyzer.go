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

package exclusions

import (
	"reflect"

	"golang.org/x/tools/go/analysis"
)

// Analyzer provides information about which types should be excluded from further analysis
// by other passes in the zerolint toolchain.
var Analyzer = &analysis.Analyzer{ //nolint:gochecknoglobals
	Name:             "exclusions",
	Doc:              "determine type exclusions for later passes",
	URL:              "https://pkg.go.dev/fillmore-labs.com/zerolint/pkg/internal/passes/exclusions",
	Run:              func(p *analysis.Pass) (any, error) { return (*pass)(p).run() },
	RunDespiteErrors: true,

	FactTypes:  []analysis.Fact{(*excludedFact)(nil)},
	ResultType: reflect.TypeFor[exclusionsResult](),
}
