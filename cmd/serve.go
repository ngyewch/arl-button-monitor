package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	buttonmonitor "github.com/ngyewch/arl-button-monitor"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

func doServe(ctx context.Context, cmd *cli.Command) error {
	buttonMonitor, err := newButtonMonitor(ctx, cmd)
	if err != nil {
		return err
	}
	defer func(buttonMonitor buttonmonitor.ButtonMonitor) {
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
			case buttonID := <-subscription.Ch():
				slog.Info("button pressed",
					slog.String("buttonId", buttonID),
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
