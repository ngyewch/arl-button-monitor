package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v3"
)

func loadConfig[T any, PT *T](ctx context.Context, cmd *cli.Command) (PT, error) {
	configFile := cmd.String(configFileFlag.Name)

	if configFile == "" {
		return nil, fmt.Errorf("%s not specified", configFileFlag.Name)
	}

	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	yamlDecoder := yaml.NewDecoder(f, yaml.Strict())
	var cfg T
	err = yamlDecoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
