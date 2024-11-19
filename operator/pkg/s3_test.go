package operator

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	retry "github.com/yetanotherco/aligned_layer/core"
	"github.com/yetanotherco/aligned_layer/core/chainio"
	s3 "github.com/yetanotherco/aligned_layer/operator/pkg/s3"
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

func TestBatchersBalances(t *testing.T) {
	cmd, err := RunStorage()
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	// To Simulate Retrieving information from S3 we create a mock http server.
	expected := "dummy data"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, expected)
	}))
	defer svr.Close()
	if err != nil {
		return
	}
	senderAddress := common.HexToAddress("0x0")

	batcher_func := s3.RequestBatch(avsWriter, &bind.CallOpts{}, senderAddress)
	_, err = batcher_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}

	batcher_func = chainio.BatcherBalances(avsWriter, &bind.CallOpts{}, senderAddress)
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

	cmd, _, err = retry_test.SetupAnvil(8545)
	if err != nil {
		t.Errorf("Error setting up Anvil: %s\n", err)
	}

	batcher_func = chainio.BatcherBalances(avsWriter, &bind.CallOpts{}, senderAddress)
	_, err = batcher_func()
	assert.Nil(t, err)

	if err := cmd.Process.Kill(); err != nil {
		t.Errorf("Error killing process: %v\n", err)
		return
	}
}
