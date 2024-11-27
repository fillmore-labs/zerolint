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
	"bufio"
	"fmt"
	"io/fs"

	"fillmore-labs.com/zerolint/pkg/set"
)

// ReadExcludes reads zero-sized types excluded from analysis from a file and returns them as a set.
func ReadExcludes(fsys fs.FS, name string) (set.Set[string], error) {
	excludes := set.New[string]()

	if name == "" {
		return excludes, nil
	}

	file, err := fsys.Open(name)
	if err != nil {
		return nil, fmt.Errorf("can't open excludes from %q: %w", name, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		expr := scanner.Bytes()
		if len(expr) == 0 || expr[0] == '#' {
			continue
		}
		excludes.Insert(string(expr))
	}
	if err2 := scanner.Err(); err2 != nil {
		return nil, fmt.Errorf("error scanning %q: %w", name, err2)
	}

	return excludes, nil
}
