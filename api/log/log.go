// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"log"

	"github.com/globocom/glbgelf"
)

// Logger implements the logger interface.
// By calling InitLog, it is initialized as a glbgelf.Logger. If one wants to change that
// and log differently (say, JSON logging for their huskyCI execution) it can be replaced
// very easily by implementing the logger interface.
var Logger logger

type logger interface {
	SendLog(extra map[string]interface{}, loglevel string, messages ...interface{}) error
}

// InitLog starts glbgelf logging.
func InitLog(developmentEnv bool, address, protocol, appName, tag string) {
	glbgelf.InitLogger(address, appName, tag, developmentEnv, protocol)
	Logger = glbgelf.Logger
}

// Info sends an info type log using glbgelf.
func Info(action, info string, msgCode int, message ...interface{}) {
	if err := Logger.SendLog(map[string]interface{}{
		"action": action,
		"info":   info},
		"INFO", MsgCode[msgCode], message); err != nil {
		ErrorGlbgelf(err)
	}
}

// Warning sends a warning type log using glbgelf.
func Warning(action, info string, msgCode int, message ...interface{}) {
	if err := Logger.SendLog(map[string]interface{}{
		"action": action,
		"info":   info},
		"WARNING", MsgCode[msgCode], message); err != nil {
		ErrorGlbgelf(err)
	}
}

// Error sends an error type log using glbgelf.
func Error(action, info string, msgCode int, message ...interface{}) {
	if err := Logger.SendLog(map[string]interface{}{
		"action": action,
		"info":   info},
		"ERROR", MsgCode[msgCode], message); err != nil {
		ErrorGlbgelf(err)
	}
}

// ErrorGlbgelf handles glbgelf error.
func ErrorGlbgelf(err error) {
	log.Println(err)
}
