package actions

import (
	"context"
	"log"

	"github.com/urfave/cli/v2"
	"github.com/yetanotherco/aligned_layer/core/config"
	"github.com/yetanotherco/aligned_layer/core/utils"
	operator "github.com/yetanotherco/aligned_layer/operator/pkg"
)

var StartFlags = []cli.Flag{
	config.ConfigFileFlag,
}

var StartCommand = &cli.Command{
	Name:        "start",
	Description: "CLI command to boot operator",
	Flags:       StartFlags,
	Action:      operatorMain,
}

func operatorMain(ctx *cli.Context) error {
	operatorConfigFilePath := ctx.String("config")
	operatorConfig := config.NewOperatorConfig(operatorConfigFilePath)
	err := utils.ReadYamlConfig(operatorConfigFilePath, &operatorConfig)
	if err != nil {
		return err
	}

	operator, err := operator.NewOperatorFromConfig(*operatorConfig)
	if err != nil {
		return err
	}

	err = operator.SendTelemetryData(ctx)
	if err != nil {
		return err
	}

	operator.Logger.Info("Operator starting...")
	err = operator.Start(context.Background())
	if err != nil {
		return err
	}

	log.Println("Operator started")

	return nil
}
