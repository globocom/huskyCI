package log

import (
	"log"
	"os"
	"strings"

	"github.com/globocom/glbgelf"
)

// InitLog starts glbgelf logging.
func InitLog() {

	isDev := true
	graylogAddr := os.Getenv("HUSKYCI_GRAYLOG_ADDR")
	gralogProto := os.Getenv("HUSKYCI_GRAYLOG_PROTO")
	appName := os.Getenv("HUSKYCI_APP_NAME")
	tags := os.Getenv("HUSKYCI_TAGS")

	if strings.EqualFold(os.Getenv("HUSKYCI_DEV"), "false") {
		isDev = false
	}

	glbgelf.InitLogger(graylogAddr, appName, tags, isDev, gralogProto)
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
