package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	buttonmonitor "github.com/ngyewch/arl-button-monitor"
	"github.com/ngyewch/arl-button-monitor/ft232h"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

func doServe(ctx context.Context, cmd *cli.Command) error {
	cfg, err := loadConfig(ctx, cmd)
	if err != nil {
		return err
	}

	buttonMonitor, err := ft232h.New(*cfg)
	if err != nil {
		return err
	}
	defer func(buttonMonitor *ft232h.ButtonMonitor) {
		_ = buttonMonitor.Close()
	}(buttonMonitor)

	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	errGroup, errGroupCtx := errgroup.WithContext(signalCtx)

	errGroup.Go(func() error {
		subscription, err := buttonMonitor.Subscribe()
		if err != nil {
			return err
		}
		defer func(subscription buttonmonitor.Subscription[string]) {
			_ = subscription.Close()
		}(subscription)

		for {
			select {
			case <-errGroupCtx.Done():
				return nil
			case err := <-subscription.ErrCh():
				return err
			case buttonId := <-subscription.Ch():
				slog.Info("button pressed",
					slog.String("buttonId", buttonId),
				)
			}
		}
	})

	err = errGroup.Wait()
	if err != nil {
		return err
	}

	return nil
}
