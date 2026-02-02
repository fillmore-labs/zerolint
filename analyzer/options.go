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

package analyzer

import (
	"log"
	"log/slog"
	"regexp"

	"fillmore-labs.com/zerolint/analyzer/level"
	"fillmore-labs.com/zerolint/internal/options"
	"fillmore-labs.com/zerolint/internal/run"
	"fillmore-labs.com/zerolint/internal/set"
)

// Option configures specific behavior of the zerolint [analysis.Analyzer].
type Option interface {
	Apply(opts *run.Options)
	LogAttr() slog.Attr
}

// Join creates a new Option joining the provided Option values.
//
// The result implements [slog.LogValuer], so the following is evaluated lazily:
//
//	slog.LogAttrs(ctx, slog.LevelInfo, "settings", Join(opts...).LogAttr())
func Join(opts ...Option) Option {
	return options.Join(opts)
}

// WithLevel is an [Option] to configure full linting.
func WithLevel(level level.LintLevel) Option {
	return levelOption{level: level}
}

type levelOption struct {
	level level.LintLevel
}

func (o levelOption) Apply(opts *run.Options) {
	opts.Level = o.level
}

// LogValue implements the [slog.LogValuer] interface.
func (o levelOption) LogAttr() slog.Attr {
	return slog.String("level", o.level.String())
}

// WithExcludes is an [Option] to configure the excluded types.
func WithExcludes(excludes []string) Option {
	return excludesOption{excludes: excludes}
}

type excludesOption struct {
	excludes []string
}

func (o excludesOption) Apply(opts *run.Options) {
	if opts.Excludes == nil {
		opts.Excludes = set.New[string]()
	}

	for _, exclude := range o.excludes {
		opts.Excludes.Add(exclude)
	}
}

// LogValue implements the [slog.LogValuer] interface.
func (o excludesOption) LogAttr() slog.Attr {
	return slog.Any("excludes", o.excludes)
}

// WithZeroTrace is an [Option] to configure tracing of zero-sized types.
func WithZeroTrace(zeroTrace bool) Option {
	return zeroTraceOption{zeroTrace: zeroTrace}
}

type zeroTraceOption struct {
	zeroTrace bool
}

func (o zeroTraceOption) Apply(opts *run.Options) {
	opts.ZeroTrace = o.zeroTrace
}

// LogValue implements the [slog.LogValuer] interface.
func (o zeroTraceOption) LogAttr() slog.Attr {
	return slog.Bool("zeroTrace", o.zeroTrace)
}

// WithGenerated is an [Option] to configure linting of generated files.
func WithGenerated(generated bool) Option {
	return generatedOption{generated: generated}
}

type generatedOption struct {
	generated bool
}

func (o generatedOption) Apply(opts *run.Options) {
	opts.Generated = o.generated
}

// LogValue implements the [slog.LogValuer] interface.
func (o generatedOption) LogAttr() slog.Attr {
	return slog.Bool("generated", o.generated)
}

// WithRegex is an [Option] to configure detecting only matching types.
func WithRegex(re *regexp.Regexp) Option {
	return reOption{re: re}
}

type reOption struct {
	re *regexp.Regexp
}

func (o reOption) Apply(opts *run.Options) {
	opts.Regex = o.re
}

// LogValue implements the [slog.LogValuer] interface.
func (o reOption) LogAttr() slog.Attr {
	re := ""
	if o.re != nil {
		re = o.re.String()
	}

	return slog.String("regex", re)
}

// WithLogger is an [Option] to configure the used logger.
func WithLogger(logger *log.Logger) Option {
	return loggerOption{logger: logger}
}

func (o loggerOption) Apply(opts *run.Options) {
	opts.Logger = o.logger
}

// LogValue implements the [slog.LogValuer] interface.
func (o loggerOption) LogAttr() slog.Attr {
	prefix := "<nil>"
	if o.logger != nil {
		prefix = o.logger.Prefix()
	}

	return slog.String("logger", prefix)
}

type loggerOption struct {
	logger *log.Logger
}

// WithExcludeComments is an [Option] to configure parsing of `zerolint:exclude` comments.
func WithExcludeComments(excludeComments bool) Option {
	return excludeCommentsOption{excludeComments: excludeComments}
}

type excludeCommentsOption struct {
	excludeComments bool
}

func (o excludeCommentsOption) Apply(opts *run.Options) {
	opts.ExcludeComments = o.excludeComments
}

// LogValue implements the [slog.LogValuer] interface.
func (o excludeCommentsOption) LogAttr() slog.Attr {
	return slog.Bool("exclude-comments", o.excludeComments)
}

// WithFlags is an [Option] to configure parsing of command-line flags.
// When enabled, command-line flags (e.g., -level, -excluded) will be parsed
// and will override any corresponding options set programmatically via other `With...` functions.
func WithFlags(flags bool) Option {
	return flagsOption{flags: flags}
}

type flagsOption struct {
	flags bool
}

func (o flagsOption) Apply(opts *run.Options) {
	opts.WithFlags = o.flags
}

// LogValue implements the [slog.LogValuer] interface.
func (o flagsOption) LogAttr() slog.Attr {
	return slog.Bool("flags", o.flags)
}
