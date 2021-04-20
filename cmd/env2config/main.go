package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/johnstarich/env2config"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	_ "github.com/johnstarich/env2config/formats"
)

type App struct {
	Configs []string
}

func main() {
	err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "e2c"))
		os.Exit(1)
	}
}

func run(args []string) error {
	var app App
	err := envconfig.Process("E2C", &app)
	if err != nil {
		return err
	}
	err = writeConfigs(app)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func writeConfigs(app App) error {
	for _, configName := range app.Configs {
		config, err := env2config.New(configName)
		if err != nil {
			return err
		}
		err = config.Write()
		if err != nil {
			return errors.Wrap(err, config.Name)
		}
	}
	return nil
}
