package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/huh/spinner"
	"github.com/urfave/cli/v2"
	"github.com/wizenheimer/cascade/internal/config"
	"github.com/wizenheimer/cascade/internal/parser"
	k8x "github.com/wizenheimer/cascade/service/kubernetes"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func main() {
	// Create a logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Name:  "cascade",
		Usage: "A CLI for managing chaos experiments",
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "Manage the server",
				Subcommands: []*cli.Command{
					{
						Name:  "start",
						Usage: "Start the server",
						Action: func(c *cli.Context) error {
							cmd := exec.Command("docker-compose", "up", "-d", "--build")
							cmd.Stdout = os.Stdout
							cmd.Stderr = os.Stderr
							return cmd.Run()
						},
					},
					{
						Name:  "stop",
						Usage: "Stop the server",
						Action: func(c *cli.Context) error {
							cmd := exec.Command("docker-compose", "down", "-v")
							cmd.Stdout = os.Stdout
							cmd.Stderr = os.Stderr
							return cmd.Run()
						},
					},
				},
			},
			{
				Name:  "create",
				Usage: "Create a new chaos experiment scenario",
				Action: func(c *cli.Context) error {
					return create(logger)
				},
			},
			{
				Name:  "exec",
				Usage: "Trigger a chaos experiment session",
				Action: func(c *cli.Context) error {
					ctx, cancel := context.WithCancel(c.Context)
					defer cancel()
					return executor(logger, ctx)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err.Error())
	}
}

func create(logger *zap.Logger) error {
	var cfg config.Config

	var outputFolder string
	var fileName string
	var createFolder bool

	form := createScenarioForm(&cfg, &outputFolder, &fileName, &createFolder)
	if err := form.Run(); err != nil {
		logger.Error(err.Error())
		return err
	}

	if err := spinner.New().
		Title("Preparing YAML ...").
		Action(func() {
			compose(cfg, fileName, outputFolder, createFolder)
		}).
		Run(); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func executor(logger *zap.Logger, ctx context.Context) error {
	var inputPath string

	form := createSessionForm(&inputPath)
	if err := form.Run(); err != nil {
		logger.Error(err.Error())
		return err
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Create a Scenario using YAML
	var config config.Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	tc, err := parser.ParseTargetConfig(&config)
	if err != nil {
		return err
	}

	rc, err := parser.ParseRuntimeConfig(&config)
	if err != nil {
		return err
	}

	cc, err := parser.ParseClusterConfig(&config)
	if err != nil {
		return err
	}

	// Create executor
	executor, err := k8x.CreateExecutor(cc, tc, rc, logger)
	if err != nil {
		return err
	}

	// Create ticker
	ticker := time.NewTicker(rc.Interval)
	defer ticker.Stop()

	// Start processing in a goroutine
	go func(executor *k8x.Executor, sessionID string, scenarioID string, ctx context.Context, next <-chan time.Time) {
		for {
			executor.Logger.Info("Chaos Session Triggered", zap.Any("Session", sessionID))
			executor.Logger.Info("Chaos Scenario", zap.Any("Scenario", scenarioID))

			// Trigger Execution
			if err = executor.Execute(ctx); err != nil {
				executor.Logger.Error(err.Error())
			}

			select {
			case <-next:
				// trigger next session
			case <-ctx.Done():
				// skip subsequent execution
				return
			}
		}
	}(executor, "nil", config.Scenario.ID, ctx, ticker.C)

	return nil
}
