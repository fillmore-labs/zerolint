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

package excludes_test

import (
	"bufio"
	"errors"
	"io/fs"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	. "fillmore-labs.com/zerolint/pkg/internal/excludes"
)

const (
	testFileGood = `# test
entry1

entry2`
	testFileEmptyContent              = ``
	testFileOnlyCommentsAndWhitespace = `
# comment1
  
	# comment2

`
	testFileWhitespaceEntries = `
  entryA  
    
  # comment
  entryB
`
)

var ErrTest = errors.New("test error")

func TestReadExcludes(t *testing.T) {
	t.Parallel()

	testfs := fstest.MapFS{
		"good.txt":                 {Data: []byte(testFileGood)},
		"empty_file.txt":           {Data: []byte(testFileEmptyContent)},
		"only_comments_blanks.txt": {Data: []byte(testFileOnlyCommentsAndWhitespace)},
		"whitespace_entries.txt":   {Data: []byte(testFileWhitespaceEntries)},
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
		{"good file", args{testfs, "good.txt"}, []string{"entry1", "entry2"}, nil},
		{"empty filename", args{testfs, ""}, nil, nil},
		{"nonexistent file", args{testfs, "nonexistent.txt"}, nil, fs.ErrNotExist},
		{"read error", args{ErrFs(ErrTest), "error_read.txt"}, nil, ErrTest},
		{"empty file", args{testfs, "empty_file.txt"}, nil, nil},
		{"only comments and blanks", args{testfs, "only_comments_blanks.txt"}, nil, nil},
		{"whitespace entries", args{testfs, "whitespace_entries.txt"}, []string{"entryA", "entryB"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ReadExcludes(tt.args.fsys, tt.args.name)
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

func TestAllText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		want      []string
		stopAfter int // Number of items to take from iterator, -1 for all
	}{
		{
			name:      "empty input",
			input:     "",
			want:      nil,
			stopAfter: -1,
		},
		{
			name:      "single line",
			input:     "hello",
			want:      []string{"hello"},
			stopAfter: -1,
		},
		{
			name:      "multiple lines",
			input:     "hello\nworld\nGo",
			want:      []string{"hello", "world", "Go"},
			stopAfter: -1,
		},
		{
			name:      "stop early",
			input:     "line1\nline2\nline3",
			want:      []string{"line1", "line2"},
			stopAfter: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			stopAfter := tt.stopAfter

			var got []string
			for s := range AllText(scanner) {
				got = append(got, s)

				if stopAfter--; stopAfter == 0 {
					break
				}
			}

			if err := scanner.Err(); err != nil {
				t.Fatalf("Scanner error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllText() got = %v, want %v", got, tt.want)
			}
		})
	}
}
