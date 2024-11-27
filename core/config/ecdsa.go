package config

import (
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"
	"os"

	ecdsa2 "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	"github.com/Layr-Labs/eigensdk-go/signer"
	"github.com/yetanotherco/aligned_layer/core/utils"
)

type EcdsaConfig struct {
	PrivateKey *ecdsa.PrivateKey
	Signer     signer.Signer
}

type EcdsaConfigFromYaml struct {
	Ecdsa struct {
		PrivateKeyStorePath     string `yaml:"private_key_store_path"`
		PrivateKeyStorePassword string `yaml:"private_key_store_password"`
	} `yaml:"ecdsa"`
}

func NewEcdsaConfig(ecdsaConfigFilePath string, chainId *big.Int) *EcdsaConfig {
	if _, err := os.Stat(ecdsaConfigFilePath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("Setup ecdsa config file does not exist")
	}

	var ecdsaConfigFromYaml EcdsaConfigFromYaml
	err := utils.ReadYamlConfig(ecdsaConfigFilePath, &ecdsaConfigFromYaml)
	if err != nil {
		log.Fatal("Error reading ecdsa config: ", err)
	}

	if ecdsaConfigFromYaml.Ecdsa.PrivateKeyStorePath == "" {
		log.Fatal("Ecdsa private key store path is empty")
	}

	ecdsaKeyPair, err := ecdsa2.ReadKey(ecdsaConfigFromYaml.Ecdsa.PrivateKeyStorePath, ecdsaConfigFromYaml.Ecdsa.PrivateKeyStorePassword)
	if err != nil {
		log.Fatal("Error reading ecdsa private key from file: ", err)
	}

	privateKeySigner, err := signer.NewPrivateKeySigner(ecdsaKeyPair, chainId)
	if err != nil {
		log.Fatal("Error creating private key signer: ", err)
	}

	return &EcdsaConfig{
		PrivateKey: ecdsaKeyPair,
		Signer:     privateKeySigner,
	}
}
