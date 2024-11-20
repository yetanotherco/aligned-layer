package operator

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	retry "github.com/yetanotherco/aligned_layer/core"
)

// EchoHandler reads and writes the request body
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, "Error writing body", http.StatusBadRequest)
		return
	}
}

// Note the httptest API requires that for restarting a server its url is empty. Given ours the default address is in use.
// Its abstracting creation of a task server works around these issues.
// NOTE: httptest panic's if the url of the server starts without being set to "" ref: https://cs.opensource.google/go/go/+/refs/tags/go1.23.3:src/net/http/httptest/server.go;l=127
func CreateTestServer() *httptest.Server {
	// To Simulate Retrieving information from S3 we create a mock http server.
	handler := http.HandlerFunc(EchoHandler) // create a listener with the desired port.
	l, err := net.Listen("tcp", "127.0.0.1:7878")
	if err != nil {
		log.Fatal(err)
	}

	svr := httptest.NewUnstartedServer(handler)

	// NewUnstartedServer creates a listener. Close that listener and replace
	// with the one we created.
	svr.Listener.Close()
	svr.Listener = l

	return svr
}

func TestRequestBatch(t *testing.T) {
	svr := CreateTestServer()

	// Start the server.
	svr.Start()

	req, err := http.NewRequestWithContext(context.Background(), "GET", svr.URL, nil)
	if err != nil {
		t.Errorf("Error creating req: %s\n", err)
	}

	batcher_func := RequestBatch(req, context.Background())
	_, err = batcher_func()
	assert.Nil(t, err)

	svr.Close()

	batcher_func = RequestBatch(req, context.Background())
	_, err = batcher_func()
	assert.NotNil(t, err)
	if _, ok := err.(retry.PermanentError); ok {
		t.Errorf("BatchersBalances Emitted non-Transient error: %s\n", err)
		return
	}
	if !strings.Contains(err.Error(), "connect: connection refused") {
		t.Errorf("BatchersBalances did not return expected error: %s\n", err)
		return
	}

	svr = CreateTestServer()
	svr.Start()
	defer svr.Close()

	batcher_func = RequestBatch(req, context.Background())
	_, err = batcher_func()
	assert.Nil(t, err)
}
