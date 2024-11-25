package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	bls "github.com/Layr-Labs/eigensdk-go/crypto/bls"
)

func main() {
	signatureArg := flag.String("signature", "", "Hex-encoded BLS signature (G1 point)")
	publicKeyArg := flag.String("publickey", "", "Hex-encoded BLS public key (G2 point)")
	messageArg := flag.String("message", "", "Hex-encoded message")

	flag.Parse()

	if *signatureArg == "" || *publicKeyArg == "" || *messageArg == "" {
		log.Fatalf("All arguments (signature, publickey, and messagehash) are required")
	}

	// Convert from hex to bytes
	signature, err := hex.DecodeString(*signatureArg)
	if err != nil {
		log.Fatalf("Failed to decode signature: %v", err)
	}

	publicKey, err := hex.DecodeString(*publicKeyArg)
	if err != nil {
		log.Fatalf("Failed to decode public key: %v", err)
	}

	messageHash, err := hex.DecodeString(*messageArg)
	if err != nil {
		log.Fatalf("Failed to decode message hash: %v", err)
	}

	isValid, err := verifySignature(signature, publicKey, messageHash)
	if err != nil {
		log.Fatalf("Error during verification: %v", err)
	}

	if isValid {
		fmt.Println("valid")
		os.Exit(0)
	} else {
		fmt.Println("invalid")
		os.Exit(1)
	}
}

func verifySignature(signature []byte, publicKey []byte, message []byte) (bool, error) {
	var messageBytes [32]byte
	copy(messageBytes[:], message)

	sign := bls.NewZeroSignature()
	_, err := sign.SetBytes(signature)
	if err != nil {
		return false, err
	}

	pubkey := bls.NewZeroG2Point()
	_, err = pubkey.SetBytes(publicKey)
	if err != nil {
		return false, err
	}

	return sign.Verify(pubkey, messageBytes)
}
