package main

import (
	"encoding/hex"
	"flag"
	"log"
	"os"

	bls "github.com/Layr-Labs/eigensdk-go/crypto/bls"
)

func main() {
	signatureArg := flag.String("signature", "", "BLS signature bytes")
	publicKeyG1X := flag.String("public-key-g1-x", "", "BLS public key on g1 affine x coord")
	publicKeyG1Y := flag.String("public-key-g1-y", "", "BLS public key on g1 affine y coord")
	publicKeyG2Arg := flag.String("public-key-g2", "", "BLS public key on g2")
	messageArg := flag.String("message", "", "Hex-encoded message")

	flag.Parse()

	if *signatureArg == "" || *publicKeyG1X == "" || *publicKeyG1Y == "" || *publicKeyG2Arg == "" || *messageArg == "" {
		log.Fatalf("All arguments (signature, publickey g1 hash, publickey g2, and messagehash) are required")
	}

	signature, err := hex.DecodeString(*signatureArg)
	if err != nil {
		log.Fatalf("Failed to decode signature: %v", err)
	}

	var pubkeyG1PointsBytes [2][]byte
	xBytes, err := hex.DecodeString(*publicKeyG1X)
	if err != nil {
		log.Fatalf("Failed to decode G1 X: %v", err)
	}
	yBytes, err := hex.DecodeString(*publicKeyG1Y)
	if err != nil {
		log.Fatalf("Failed to decode G1 Y: %v", err)
	}
	pubkeyG1PointsBytes[0] = xBytes
	pubkeyG1PointsBytes[1] = yBytes

	pubkeyG2Bytes, err := hex.DecodeString(*publicKeyG2Arg)
	if err != nil {
		log.Fatalf("Failed to decode pubkey: %v", err)
	}

	messageHash, err := hex.DecodeString(*messageArg)
	if err != nil {
		log.Fatalf("Failed to decode message hash: %v", err)
	}

	isValid, err := verifySignature(signature, pubkeyG1PointsBytes, pubkeyG2Bytes, messageHash)
	if err != nil {
		log.Fatalf("Error during verification: %v", err)
	}

	if isValid {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func verifySignature(signature []byte, pubkeyG1PointsBytes [2][]byte, pubkeyG2Bytes []byte, message []byte) (bool, error) {
	pubkeyG1 := bls.NewZeroG1Point()
	pubkeyG1.X.SetBytes(pubkeyG1PointsBytes[0])
	pubkeyG1.Y.SetBytes(pubkeyG1PointsBytes[1])

	pubkeyG2 := bls.NewZeroG2Point()
	_, err := pubkeyG2.SetBytes(pubkeyG2Bytes)
	if err != nil {
		return false, err
	}

	var messageBytes [32]byte
	copy(messageBytes[:], message[:])

	sign := bls.NewZeroSignature()
	_, err = sign.SetBytes(signature)
	if err != nil {
		return false, err
	}

	// verify the equivalence between the points in the generators
	valid, err := pubkeyG1.VerifyEquivalence(pubkeyG2)
	if err != nil || !valid {
		return false, err
	}

	return sign.Verify(pubkeyG2, messageBytes)
}
