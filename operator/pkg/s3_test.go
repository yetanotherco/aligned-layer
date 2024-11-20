package operator

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	retry "github.com/yetanotherco/aligned_layer/core"
)

// Function wrapper around `make run_storage`
func RunStorage() (*exec.Cmd, error) {

	// Create a command
	cmd := exec.Command("make", "run_storage")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Run the command
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	// Delay needed for anvil to start
	time.Sleep(750 * time.Millisecond)

	return cmd, nil
}
func TestRequestBatch(t *testing.T) {
	/*
		cmd, err := RunStorage()
		if err != nil {
			t.Errorf("Error setting up Anvil: %s\n", err)
		}
	*/
	// To Simulate Retrieving information from S3 we create a mock http server.
	expected := "dummy data"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("error %v %v", w, expected)
	}))
	defer svr.Close()

	req, err := http.NewRequestWithContext(context.Background(), "GET", svr.URL, nil)
	if err != nil {
		t.Errorf("Error creating req: %s\n", err)
	}

	batcher_func := RequestBatch(req, context.Background())
	_, err = batcher_func()
	assert.Nil(t, err)

	svr.Close()
	/*
		if err := cmd.Process.Kill(); err != nil {
			t.Errorf("Error killing process: %v\n", err)
			return
		}
	*/

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

	svr.Start()
	/*
		cmd, err = RunStorage()
		if err != nil {
			t.Errorf("Error setting up Anvil: %s\n", err)
		}
	*/

	batcher_func = RequestBatch(req, context.Background())
	_, err = batcher_func()
	assert.Nil(t, err)

	/*
		if err := cmd.Process.Kill(); err != nil {
			t.Errorf("Error killing process: %v\n", err)
			return
		}
	*/
}
