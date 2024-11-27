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

package excludes_test

import (
	"errors"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"

	"fillmore-labs.com/zerolint/pkg/internal/excludes"
)

const testFile = `# test
entry1

entry2`

var ErrTest = errors.New("test error")

func TestReadExcludes(t *testing.T) {
	t.Parallel()

	testfs := fstest.MapFS{
		"test": {Data: []byte(testFile)},
	}

	type args struct {
		fsys fs.FS
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr error
	}{
		{"test", args{testfs, "test"}, []string{"entry1", "entry2"}, nil},
		{"empty", args{testfs, ""}, nil, nil},
		{"nonexistent", args{testfs, "nonexistent"}, nil, fs.ErrNotExist},
		{"invalidRead", args{ErrFs(ErrTest), "err"}, nil, ErrTest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := excludes.ReadExcludes(tt.args.fsys, tt.args.name)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ReadExcludes() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadExcludes() = %v, want %v", got, tt.want)
			}
		})
	}
}
