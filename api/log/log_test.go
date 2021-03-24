package log_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	apiContext "github.com/globocom/huskyCI/api/context"

	"github.com/globocom/huskyCI/api/log"
)

func TestInitLog(t *testing.T) {
	apiContext.APIConfiguration = &apiContext.APIConfig{
		GraylogConfig: &apiContext.GraylogConfig{
			DevelopmentEnv: true,
			AppName:        "log_test",
			Tag:            "log_test",
		},
	}

	log.InitLog(true, "", "", "log_test", "log_test")

	if log.Logger == nil {
		t.Error("expected logger to be initialized, but it wasn't")
		return
	}
}

func TestLog(t *testing.T) {

	testCases := []struct {
		logger interface {
			SendLog(extra map[string]interface{}, loglevel string, messages ...interface{}) error
		}
		name         string
		wantAction   string
		wantInfo     string
		wantMsgCode  int
		wantLogLevel string
		wantMessages []string
		logFunc      func(action, info string, msgCode int, message ...interface{})
		err          error
		wantErr      string
	}{
		{
			name:         "Testing log.Info",
			logger:       &stubLogger{},
			wantAction:   "action",
			wantInfo:     "info",
			wantMsgCode:  11,
			wantLogLevel: "INFO",
			wantMessages: []string{"got some info!"},
			logFunc:      log.Info,
		},
		{
			name:         "Testing log.Error",
			logger:       &stubLogger{},
			wantAction:   "action",
			wantInfo:     "err",
			wantMsgCode:  11,
			wantLogLevel: "ERROR",
			wantMessages: []string{"got some error!"},
			logFunc:      log.Error,
		},
		{
			name:         "Testing log.Error fail",
			logger:       &stubLogger{err: errors.New("server is down!")},
			wantAction:   "action",
			wantInfo:     "err",
			wantMsgCode:  11,
			wantLogLevel: "ERROR",
			wantMessages: []string{"got some error!"},
			logFunc:      log.Error,
			wantErr:      "server is down!",
		},
		{
			name:         "Testing log.Warning",
			logger:       &stubLogger{},
			wantAction:   "action",
			wantInfo:     "err",
			wantMsgCode:  11,
			wantLogLevel: "WARNING",
			wantMessages: []string{"got some warning!"},
			logFunc:      log.Warning,
		},
		{
			name:         "Testing log.Warning fail",
			logger:       &stubLogger{err: errors.New("server is down!")},
			wantAction:   "action",
			wantInfo:     "err",
			wantMsgCode:  11,
			wantLogLevel: "WARNING",
			wantMessages: []string{"got some warning!"},
			logFunc:      log.Warning,
			wantErr:      "server is down!",
		},
		{
			name:         "Testing 'log server is down!'",
			logger:       &stubLogger{},
			wantAction:   "action",
			wantInfo:     "info",
			wantMsgCode:  11,
			wantLogLevel: "INFO",
			wantMessages: []string{"got some info!"},
			logFunc:      log.Info,
			err:          errors.New("server is down!"),
			wantErr:      "server is down!",
		},
		{
			name:         "Testing Send errored out",
			logger:       &stubLogger{err: errors.New("server is down!")},
			wantAction:   "action",
			wantInfo:     "info",
			wantMsgCode:  11,
			wantLogLevel: "INFO",
			wantMessages: []string{"got some info!"},
			logFunc:      log.Info,
			wantErr:      "server is down!",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			log.Logger = testCase.logger
			testCase.logFunc(testCase.wantAction, testCase.wantInfo, testCase.wantMsgCode, testCase.wantMessages)

			stub := testCase.logger.(*stubLogger)
			extra := stub.calledWith["extra"].(map[string]interface{})
			if got, ok := extra["action"]; !ok || got != testCase.wantAction {
				t.Errorf("in action key, we expected %s; but got %s", testCase.wantAction, got)
				return
			}
			if got, ok := extra["info"]; !ok || got != testCase.wantInfo {
				t.Errorf("in info key, we expected %s; but got %s", testCase.wantInfo, got)
				return
			}
			gotLogLevel := stub.calledWith["loglevel"]
			if gotLogLevel != testCase.wantLogLevel {
				t.Errorf("in loglevel, we expected %s; but got %s", testCase.wantLogLevel, gotLogLevel)
				return
			}
			gotMessages := stub.calledWith["messages"]
			if reflect.DeepEqual(gotMessages, testCase.wantMessages) {
				t.Errorf("in messages, we expected %s; but got %s", strings.Join(testCase.wantMessages, " "), strings.Join(gotMessages.([]string), " "))
				return
			}

			var err error
			if testCase.err != nil {
				err = testCase.err
			} else if stub.err != nil {
				err = stub.err
			}
			if err != nil && err.Error() != testCase.wantErr {
				t.Errorf("in err case, we expected %s; but got %s", testCase.err.Error(), testCase.wantErr)
				return
			}

		})
	}

}

type stubLogger struct {
	calledWith map[string]interface{}
	err        error
}

func (s *stubLogger) SendLog(extra map[string]interface{}, loglevel string, messages ...interface{}) error {
	s.calledWith = map[string]interface{}{"extra": extra, "loglevel": loglevel, "messages": messages}
	return s.err
}
