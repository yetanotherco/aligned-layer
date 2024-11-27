package actions

import (
	"context"
	operator "github.com/yetanotherco/aligned_layer/operator/pkg"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"
	"github.com/yetanotherco/aligned_layer/core/config"
)

var registerFlags = []cli.Flag{
	config.ConfigFileFlag,
}

var RegisterCommand = &cli.Command{
	Name:        "register",
	Usage:       "Register operator with Aligned Layer",
	Description: "CLI command to register opeartor with Aligned Layer",
	Flags:       registerFlags,
	Action:      registerOperatorMain,
}

func registerOperatorMain(ctx *cli.Context) error {
	operatorConfig := config.NewOperatorConfig(ctx.String(config.ConfigFileFlag.Name))
	ecdsaConfig := config.NewEcdsaConfig(ctx.String(config.ConfigFileFlag.Name), operatorConfig.BaseConfig.ChainId)

	quorumNumbers := []byte{0}

	// Generate salt and expiry
	privateKeyBytes := []byte(operatorConfig.BlsConfig.KeyPair.PrivKey.String())
	salt := [32]byte{}

	copy(salt[:], crypto.Keccak256([]byte("churn"), []byte(time.Now().String()), quorumNumbers, privateKeyBytes))

	err := operator.RegisterOperator(context.Background(), operatorConfig, ecdsaConfig, salt)
	if err != nil {
		operatorConfig.BaseConfig.Logger.Error("Failed to register operator", "err", err)
		return err
	}

	return nil
}
