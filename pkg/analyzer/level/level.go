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

package level

import (
	"errors"
	"strings"
)

// LintLevel represents the severity level of a linting issue.
type LintLevel int

const (
	// Default is the standard linting level that applies basic checks to the code.
	// It is the default level used when no specific linting level is specified.
	Default LintLevel = iota

	// Extended finds more usages of pointers to zero-sized types.
	Extended

	// Full finds most usages of pointers to zero-sized types.
	// This level probably need an excludes file.
	Full
)

// ErrUnknownLintLevel is an error returned when an unrecognized or invalid
// lint level is encountered.
var ErrUnknownLintLevel = errors.New("unknown lint level")

// UnmarshalText implements encoding.TextUnmarshaler.
func (l *LintLevel) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "default", "0":
		*l = Default

		return nil

	case "extended", "1":
		*l = Extended

		return nil

	case "full", "2":
		*l = Full

		return nil

	default:
		return ErrUnknownLintLevel
	}
}

// MarshalText implements encoding.TextMarshaler.
func (l LintLevel) MarshalText() ([]byte, error) {
	switch l {
	case Default:
		return []byte("default"), nil

	case Extended:
		return []byte("extended"), nil

	case Full:
		return []byte("full"), nil

	default:
		return nil, ErrUnknownLintLevel
	}
}
