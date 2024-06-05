package merkle_tree

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestVerifyMerkleTreeBatch(t *testing.T) {
	batchFile, err := os.Open("./lib/test_files/7a3d9215cfac21a4b0e94382e53a9f26bc23ed990f9c850a31ccf3a65aec1466.json")
	if err != nil {
		t.Fatalf("Error opening batch file: %v", err)
	}

	byteValue, err := io.ReadAll(batchFile)
	if err != nil {
		t.Fatalf("Error reading batch file: %v", err)
	}

	hexMerkleRootStr := "66f2b058f5eaceff958a1feff5edf225be9a422c1121bc7ba4fceea2549fac02"

	byteSliceFromMerkleRoot, err := hex.DecodeString(hexMerkleRootStr)
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
		return
	}

	var merkleRoot [32]byte
	copy(merkleRoot[:], byteSliceFromMerkleRoot)

	if !VerifyMerkleTreeBatch(([MaxBatchSize]byte)(byteValue), uint32(len(byteValue)), merkleRoot) {
		t.Errorf("Batch did not verify")
	}

}
