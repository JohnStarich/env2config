package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	var configErrs []string
	for _, configName := range app.Configs {
		err := writeConfig(configName)
		if err != nil {
			configErrs = append(configErrs, err.Error())
		}
	}
	if len(configErrs) > 0 {
		return errors.New("Failed to generate configs:\n\n" + strings.Join(configErrs, "\n\n"))
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

func writeConfig(name string) error {
	config, err := env2config.New(name)
	if err != nil {
		return err
	}
	err = config.Write()
	return errors.Wrap(err, config.Name)
}
