// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"log"

	"github.com/globocom/glbgelf"
	apiContext "github.com/globocom/huskyCI/api/context"
)

// InitLog starts glbgelf logging.
func InitLog() {
	graylogConfig := apiContext.APIConfiguration.GraylogConfig
	isDev := graylogConfig.DevelopmentEnv
	graylogAddr := graylogConfig.Address
	gralogProto := graylogConfig.Protocol
	appName := graylogConfig.AppName
	tag := graylogConfig.Tag
	glbgelf.InitLogger(graylogAddr, appName, tag, isDev, gralogProto)
}

// Info sends an info type log using glbgelf.
func Info(action, info string, msgCode int, message ...interface{}) {
	if err := glbgelf.Logger.SendLog(map[string]interface{}{
		"action": action,
		"info":   info},
		"INFO", MsgCode[msgCode], message); err != nil {
		ErrorGlbgelf(err)
	}
}

// Warning sends a warning type log using glbgelf.
func Warning(action, info string, msgCode int, message ...interface{}) {
	if err := glbgelf.Logger.SendLog(map[string]interface{}{
		"action": action,
		"info":   info},
		"WARNING", MsgCode[msgCode], message); err != nil {
		ErrorGlbgelf(err)
	}
}

// Error sends an error type log using glbgelf.
func Error(action, info string, msgCode int, message ...interface{}) {
	if err := glbgelf.Logger.SendLog(map[string]interface{}{
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
