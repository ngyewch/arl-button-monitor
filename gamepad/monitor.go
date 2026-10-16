package gamepad

import (
	"context"
	"log/slog"

	"github.com/doingharm/go-gamepad-bus"
	buttonmonitor "github.com/ngyewch/arl-button-monitor"
	"github.com/ngyewch/arl-button-monitor/base"
)

type ButtonMonitor struct {
	cfg          Config
	base         *base.ButtonMonitor
	gamepadBus   gamepads.Bus
	gamepadErrCh <-chan error
	cancel       func()
}

func New(cfg Config) (*ButtonMonitor, error) {
	baseButtonMonitor, err := base.New()
	if err != nil {
		return nil, err
	}

	gamepadBus, gamepadErrCh, err := gamepads.New()
	if err != nil {
		_ = baseButtonMonitor.Close()
		return nil, err
	}

	runCtx, cancel := context.WithCancel(context.Background())

	bm := &ButtonMonitor{
		cfg:          cfg,
		base:         baseButtonMonitor,
		gamepadBus:   gamepadBus,
		gamepadErrCh: gamepadErrCh,
		cancel:       cancel,
	}

	go bm.run(runCtx)

	return bm, nil
}

func (bm *ButtonMonitor) Close() error {
	if bm.cancel != nil {
		bm.cancel()
	}
	_ = bm.gamepadBus.Close
	_ = bm.base.Close()
	return nil
}

func (bm *ButtonMonitor) Subscribe() (buttonmonitor.Subscription[string], error) {
	return bm.base.Subscribe()
}

func (bm *ButtonMonitor) run(ctx context.Context) {
	eventChannel := bm.gamepadBus.NewEventChannel()
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-bm.gamepadErrCh:
			if err != nil {
				slog.Error("gamepad bus error",
					slog.Any("err", err),
				)
			}
		case err := <-eventChannel.ErrCh:
			if err != nil {
				slog.Error("gamepad bus error",
					slog.Any("err", err),
				)
			}
		case event := <-eventChannel.Ch:
			switch event.Type {
			case gamepads.ConnectEventType:
				err := bm.gamepadBus.Subscribe(event.ID)
				if err != nil {
					slog.Error("could not subscribe to gamepad",
						slog.String("id", event.ID),
						slog.Any("err", err),
					)
				}
			case gamepads.DisconnectEventType:
				err := bm.gamepadBus.Unsubscribe(event.ID)
				if err != nil {
					slog.Error("could not unsubscribe from gamepad",
						slog.String("id", event.ID),
						slog.Any("err", err),
					)
				}
			case gamepads.ControlEventType:
				controlEvent := event.Data.(gamepads.ControlEvent)
				if controlEvent.Type == gamepads.Button {
					if controlEvent.Value == 1 {
						for _, buttonCfg := range bm.cfg.Buttons {
							if (event.ID == buttonCfg.GamepadID) && (controlEvent.Index == buttonCfg.ButtonID) {
								err := bm.base.Publish(buttonCfg.ID)
								if err != nil {
									slog.Error("could not publish gamepad button click event",
										slog.Any("err", err),
									)
								}
							}
						}
					}
				}
			}
		}
	}
}
