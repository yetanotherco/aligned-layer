package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/yetanotherco/aligned_layer/aggregator/internal/pkg"
	"github.com/yetanotherco/aligned_layer/core/config"
	"github.com/yetanotherco/aligned_layer/core/sched"
	"github.com/yetanotherco/aligned_layer/core/supervisor"
)

var (
	// Version is the version of the binary.
	Version   string
	GitCommit string
	GitDate   string
)

var flags = []cli.Flag{
	config.ConfigFileFlag,
}

func main() {
	app := cli.NewApp()

	app.Flags = flags
	app.Version = fmt.Sprintf("%s-%s-%s", Version, GitCommit, GitDate)
	app.Name = "aligned-layer-aggregator"
	app.Usage = "Aligned Layer Aggregator"
	app.Description = "Service that aggregates signed responses from operator nodes."
	app.Action = aggregatorMain

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln("Application failed.", "Message:", err)
	}
}

func aggregatorMain(ctx *cli.Context) error {

	configFilePath := ctx.String(config.ConfigFileFlag.Name)
	aggregatorConfig := config.NewAggregatorConfig(configFilePath)

	aggregator, err := pkg.NewAggregator(*aggregatorConfig)
	if err != nil {
		aggregatorConfig.BaseConfig.Logger.Error("Cannot create aggregator", "err", err)
		return err
	}

	gcPeriod := aggregatorConfig.Aggregator.GarbageCollectorPeriod
	aggregatorConfig.BaseConfig.Logger.Info(fmt.Sprintf("- Removing finalized Task Infos from Maps every %v", gcPeriod))
	lastIdxDeleted := uint32(0)

	// Supervisor revives garbage collector
	sched.Every(gcPeriod, func() error {
		var err error
		lastIdxDeleted, err = aggregator.ClearTasksFromMaps(lastIdxDeleted)
		return err
	})

	// Listen for new task created in the ServiceManager contract in a separate goroutine, both V1 and V2 subscriptions:
	supervisor.Serve(func() {
		listenErr := aggregator.SubscribeToNewTasks()
		if listenErr != nil {
			aggregatorConfig.BaseConfig.Logger.Fatal("Error subscribing for new tasks", "err", listenErr)
		}
	}, "subscriber")

	err = aggregator.Start(context.Background())

	return err
}
