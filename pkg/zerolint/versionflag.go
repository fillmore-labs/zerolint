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

package zerolint

import (
	"fmt"
	"os"
	"runtime/debug"
)

// versionFlag represents a [flag] to print version information and exit the program.
type versionFlag struct{}

func (versionFlag) IsBoolFlag() bool { return true }
func (versionFlag) String() string   { return "true" }
func (versionFlag) Set(_ string) error {
	const progname = Name

	if bi, ok := debug.ReadBuildInfo(); ok {
		fmt.Printf("%s version %s build with %s\n",
			progname, bi.Main.Version, bi.GoVersion)
	} else {
		fmt.Printf("%s version (unknown)\n", progname)
	}

	os.Exit(0)

	return nil
}
