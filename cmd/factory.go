package main

import (
	"context"
	"fmt"

	buttonmonitor "github.com/ngyewch/arl-button-monitor"
	"github.com/ngyewch/arl-button-monitor/dummy"
	"github.com/ngyewch/arl-button-monitor/ft232h"
	"github.com/ngyewch/arl-button-monitor/gamepad"
	"github.com/urfave/cli/v3"
)

func newButtonMonitor(ctx context.Context, cmd *cli.Command) (buttonmonitor.ButtonMonitor, error) {
	buttonMonitorType := cmd.String(typeFlag.Name)

	switch buttonMonitorType {
	case "dummy":
		return dummy.New()
	case "ft232h":
		cfg, err := loadConfig[ft232h.Config, *ft232h.Config](ctx, cmd)
		if err != nil {
			return nil, err
		}
		return ft232h.New(*cfg)
	case "gamepad":
		cfg, err := loadConfig[gamepad.Config, *gamepad.Config](ctx, cmd)
		if err != nil {
			return nil, err
		}
		return gamepad.New(*cfg)
	default:
		return nil, fmt.Errorf("unknown type: %s", buttonMonitorType)
	}
}
