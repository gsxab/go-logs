/*
 * SPDX-License-Identifier: MIT
 *
 * Copyright (c) 2023 Gsxab
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package logs

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type LogLevel int8

const (
	AllLevels LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	NoLevels
)

var (
	logLevel = InfoLevel
	errorLog = log.New(os.Stdout, "ERROR ", log.Ltime|log.Llongfile)
	warnLog  = log.New(os.Stdout, "WARN  ", log.Ltime|log.Llongfile)
	infoLog  = log.New(os.Stdout, "INFO  ", log.Ltime|log.Llongfile)
	debugLog = log.New(os.Stdout, "DEBUG ", log.Ltime|log.Llongfile)
)

func resetLoggerForLevels(from, to LogLevel, logger io.Writer) {
	switch from {
	case AllLevels:
		if to == DebugLevel {
			break
		}
		fallthrough
	case DebugLevel:
		debugLog = log.New(logger, "DEBUG ", log.Ltime|log.Llongfile)
		if to == InfoLevel {
			break
		}
		fallthrough
	case InfoLevel:
		infoLog = log.New(logger, "INFO  ", log.Ltime|log.Llongfile)
		if to == WarnLevel {
			break
		}
		fallthrough
	case WarnLevel:
		warnLog = log.New(logger, "WARN  ", log.Ltime|log.Llongfile)
		if to == ErrorLevel {
			break
		}
		fallthrough
	case ErrorLevel:
		errorLog = log.New(logger, "ERROR ", log.Ltime|log.Llongfile)
	}
}

func enableLevels(from, to LogLevel) {
	resetLoggerForLevels(from, to, os.Stdout)
}

func disableLevels(from, to LogLevel) {
	resetLoggerForLevels(from, to, io.Discard)
}

func SetLevel(level LogLevel) {
	if level == logLevel {
		return
	} else if level < logLevel {
		enableLevels(level, logLevel)
	} else {
		disableLevels(logLevel, level)
	}
	logLevel = level
}

func format(kvs ...any) string {
	sb := &strings.Builder{}
	//_, _ = fmt.Fprintf(sb, "Module=%s Function=%s Message=%s", module, function, message)
	for i := 0; i < len(kvs)-1; i += 2 {
		_, _ = fmt.Fprintf(sb, " %s=%+v", kvs[i], kvs[i+1])
	}
	return sb.String()
}

func makeKvLogFunc(level LogLevel, pLogger **log.Logger) func(ctx context.Context, kvs ...any) {
	return func(ctx context.Context, kvs ...any) {
		str := format(kvs...)
		_ = (*pLogger).Output(3, str)
	}
}

func makeBasicLogFunc(f func(context.Context, ...any)) func(ctx context.Context, kvs ...any) {
	return func(ctx context.Context, kvs ...any) {
		f(ctx, kvs...)
	}
}

func makeMsgLogFunc(f func(context.Context, ...any)) func(ctx context.Context, msg string, kvs ...any) {
	return func(ctx context.Context, msg string, kvs ...any) {
		kvs = append([]any{"msg", msg}, kvs...)
		f(ctx, kvs...)
	}
}

func makeErrLogFunc(f func(context.Context, ...any)) func(ctx context.Context, err error, kvs ...any) {
	return func(ctx context.Context, err error, kvs ...any) {
		kvs = append(kvs, "err", err.Error())
		f(ctx, kvs...)
	}
}

func makeMsgErrLogFunc(f func(context.Context, ...any)) func(ctx context.Context, msg string, err error, kvs ...any) {
	return func(ctx context.Context, msg string, err error, kvs ...any) {
		kvs = append([]any{"msg", msg}, kvs...)
		kvs = append(kvs, "err", err.Error())
		f(ctx, kvs...)
	}
}

//goland:noinspection GoUnusedGlobalVariable
var (
	errorLogger = makeKvLogFunc(ErrorLevel, &errorLog)
	warnLogger  = makeKvLogFunc(WarnLevel, &warnLog)
	infoLogger  = makeKvLogFunc(InfoLevel, &infoLog)
	debugLogger = makeKvLogFunc(DebugLevel, &debugLog)

	ErrorR = makeBasicLogFunc(errorLogger)
	WarnR  = makeBasicLogFunc(warnLogger)
	InfoR  = makeBasicLogFunc(infoLogger)
	DebugR = makeBasicLogFunc(debugLogger)

	ErrorM = makeMsgLogFunc(errorLogger)
	WarnM  = makeMsgLogFunc(warnLogger)
	InfoM  = makeMsgLogFunc(infoLogger)
	DebugM = makeMsgLogFunc(debugLogger)

	ErrorE = makeErrLogFunc(errorLogger)
	WarnE  = makeErrLogFunc(warnLogger)
	InfoE  = makeErrLogFunc(infoLogger)
	DebugE = makeErrLogFunc(debugLogger)

	Error = makeMsgErrLogFunc(errorLogger)
	Warn  = makeMsgErrLogFunc(warnLogger)
	Info  = makeMsgErrLogFunc(infoLogger)
	Debug = makeMsgErrLogFunc(debugLogger)
)
