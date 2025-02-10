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

package excludes

import (
	"bufio"
	"fmt"
	"io/fs"
	"iter"
	"strings"
)

// ReadExcludes reads zero-sized types excluded from analysis from a file and returns them as a list.
func ReadExcludes(fsys fs.FS, name string) ([]string, error) {
	if name == "" {
		return nil, nil
	}

	file, err := fsys.Open(name)
	if err != nil {
		return nil, fmt.Errorf("can't open excludes from %q: %w", name, err)
	}
	defer file.Close()

	var excludes []string //nolint:prealloc
	scanner := bufio.NewScanner(file)
	for line := range AllText(scanner) {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		excludes = append(excludes, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning %q: %w", name, err)
	}

	return excludes, nil
}

// AllText iterates over the tokens generated by a scanner.
func AllText(scanner *bufio.Scanner) iter.Seq[string] {
	return func(yield func(string) bool) {
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				break
			}
		}
	}
}
