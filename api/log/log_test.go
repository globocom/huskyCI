package log_test

import (
	"bufio"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/globocom/glbgelf"
	apiContext "github.com/globocom/huskyCI/api/context"

	"github.com/globocom/huskyCI/api/log"
)

func startGelfServer(t *testing.T, wg *sync.WaitGroup, msgCount int) {
	defer wg.Done()

	// set apiContext vars
	listener, err := net.Listen("tcp", ":12201")
	if err != nil {
		t.Errorf("expected to listening server, but failed: %s", err)
		return
	}
	// We could use ExampleInitLog, however, I think it is better if we read it from the stream.
	conn, err := listener.Accept()
	if err != nil {
		t.Errorf("could not accept messages: %s", err)
		return
	}
	var gotMsgs int
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	reader := bufio.NewReader(conn)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			t.Errorf("could not read bytes from conn: %s", err)
			return
		}
		if b == '\u0000' { // as per gelf docs: a message is finished when it is NUL terminated
			gotMsgs++
		}
		if gotMsgs == msgCount {
			break
		}
	}
	if gotMsgs != msgCount {
		t.Errorf("expected %d msgs, but got %d", msgCount, gotMsgs)
		return
	}
	conn.Close()
}

func TestLog(t *testing.T) {
	apiContext.APIConfiguration = &apiContext.APIConfig{
		GraylogConfig: &apiContext.GraylogConfig{
			DevelopmentEnv: false,
			Address:        "localhost:12201",
			Protocol:       "tcp",
			AppName:        "log_test",
			Tag:            "log_test",
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go startGelfServer(t,
		wg,
		5, // logging to stdout, starting husky, logging to localhost, info, warning, error.
	)
	log.InitLog()

	if glbgelf.Logger == nil {
		t.Error("expected logger to be initialized, but it wasn't")
		return
	}

	// Try to send a simple message to our mock server
	log.Info("action", "info", 11, "infoTest")
	log.Warning("action", "info", 101, "warnTest")
	log.Error("action", "info", 1006, "errorTest")
	wg.Wait()
}
