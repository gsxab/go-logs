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
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

type StreamType int8

const (
	StreamTypeDiscard StreamType = iota
	StreamTypeStdout
	StreamTypeStderr
	StreamTypeFileWriter
)

type ConfigItem struct {
	Level      LogLevel     `json:"level"`
	StreamType StreamType   `json:"stream_type"`
	Params     ConfigParams `json:"params,omitempty"`
}

type Config struct {
	Items []*ConfigItem
}

type ConfigParams struct {
	Filename string `json:"filename,omitempty"`
	Perm     *int32 `json:"perm,omitempty"`
}

func LoadConfig(str []byte) (*Config, error) {
	config := &Config{}
	err := json.Unmarshal(str, config)
	if err != nil {
		return nil, err
	}
	err = UseConfig(config)
	return config, err
}

func UseConfig(config *Config) error {
	var debugWriters, infoWriters, warnWriters, errorWriters []io.Writer
	for _, item := range config.Items {
		var writer io.Writer
		switch item.StreamType {
		case StreamTypeDiscard:
			writer = io.Discard
		case StreamTypeStdout:
			writer = os.Stdout
		case StreamTypeStderr:
			writer = os.Stderr
		case StreamTypeFileWriter:
			filename := item.Params.Filename
			if filename == "" {
				return errors.New("filename is required")
			}

			perm := item.Params.Perm
			permInt := int32(0x644)
			if perm != nil {
				permInt = *perm
			}

			file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.FileMode(permInt))
			if err != nil {
				return err
			}

			writer = file
		}
		switch item.Level {
		case AllLevels:
			fallthrough
		case DebugLevel:
			debugWriters = append(debugWriters, writer)
			fallthrough
		case InfoLevel:
			infoWriters = append(infoWriters, writer)
			fallthrough
		case WarnLevel:
			warnWriters = append(warnWriters, writer)
			fallthrough
		case ErrorLevel:
			errorWriters = append(errorWriters, writer)
		}
	}
	debugLog = log.New(io.MultiWriter(debugWriters...), "DEBUG ", log.Ltime|log.Llongfile)
	infoLog = log.New(io.MultiWriter(infoWriters...), "INFO  ", log.Ltime|log.Llongfile)
	warnLog = log.New(io.MultiWriter(warnWriters...), "WARN  ", log.Ltime|log.Llongfile)
	errorLog = log.New(io.MultiWriter(errorWriters...), "ERROR ", log.Ltime|log.Llongfile)
	return nil
}
