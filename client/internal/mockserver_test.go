package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/open-telemetry/opamp-go/protobufs"
)

func TestMockServerHandleMessageAfterClose(t *testing.T) {
	srv := StartMockServer(t)
	srv.EnableExpectMode()
	srv.Close()

	msgBytes, err := proto.Marshal(&protobufs.AgentToServer{InstanceUid: []byte("12345678901234567890123456")})
	require.NoError(t, err)

	// The old code panicked on a nil handler read from the closed expectedHandlers
	// channel and then deadlocked in the deferred send, so run the call in a
	// goroutine and enforce a timeout.
	done := make(chan []byte, 1)
	go func() {
		done <- srv.handleReceivedBytes(msgBytes, false)
	}()

	select {
	case response := <-done:
		assert.Nil(t, response)
	case <-time.After(5 * time.Second):
		t.Fatal("handleReceivedBytes did not return after Close()")
	}
}
