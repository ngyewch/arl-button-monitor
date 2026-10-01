package ft232h

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/F2077/go-pubsub/pubsub"
	buttonmonitor "github.com/ngyewch/arl-button-monitor"
	"github.com/zmwangx/debounce"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

const (
	topicButtonPressed = "buttonPressed"
)

type Button struct {
	id           string
	pin          gpio.PinIO
	pressedLevel gpio.Level
	handler      func()
}

type ButtonMonitor struct {
	cfg       Config
	buttons   []Button
	broker    *pubsub.Broker[string]
	publisher *pubsub.Publisher[string]
	cancel    func()
}

func New(cfg Config) (*ButtonMonitor, error) {
	ft232h := findFT232H()
	if ft232h == nil {
		return nil, fmt.Errorf("could not find FT232H device on the USB bus")
	}
	broker, err := pubsub.NewBroker[string]()
	if err != nil {
		return nil, err
	}
	publisher := pubsub.NewPublisher[string](broker)
	var buttons []Button
	for _, buttonConfig := range cfg.Buttons {
		pin := gpioreg.ByName(buttonConfig.GPIOPinName)
		if pin == nil {
			return nil, fmt.Errorf("could not find GPIO pin %s", buttonConfig.GPIOPinName)
		}
		err := pin.In(gpio.PullNoChange, gpio.NoEdge)
		if err != nil {
			return nil, err
		}
		pressedLevel, err := levelFromString(buttonConfig.PressedLevel)
		if err != nil {
			return nil, err
		}
		debounced, _ := debounce.Debounce(func() {
			slog.Debug("button pressed",
				slog.String("id", buttonConfig.ID),
			)
			err = publisher.Publish(topicButtonPressed, buttonConfig.ID)
			if err != nil {
				slog.Error("could not publish button press event",
					slog.Any("err", err),
				)
			}
		}, cfg.DebouncePeriod, debounce.WithLeading(true), debounce.WithTrailing(false))
		buttons = append(buttons, Button{
			id:           buttonConfig.ID,
			pin:          pin,
			pressedLevel: pressedLevel,
			handler:      debounced,
		})
	}
	runCtx, cancel := context.WithCancel(context.Background())
	bm := &ButtonMonitor{
		cfg:       cfg,
		buttons:   buttons,
		broker:    broker,
		publisher: publisher,
		cancel:    cancel,
	}
	go bm.run(runCtx)
	return bm, nil
}

func (bm *ButtonMonitor) Close() error {
	if bm.cancel != nil {
		bm.cancel()
	}
	return nil
}

func (bm *ButtonMonitor) Subscribe() (buttonmonitor.Subscription[string], error) {
	subscriber := pubsub.NewSubscriber[string](bm.broker)
	subscription, err := subscriber.Subscribe(topicButtonPressed, pubsub.WithChannelSize[string](pubsub.DefaultChannelSize))
	if err != nil {
		_ = subscriber.Close()
		return nil, err
	}
	return &customSubscription[string]{
		subscriber:   subscriber,
		subscription: subscription,
	}, nil
}

func (bm *ButtonMonitor) run(ctx context.Context) {
	for {
		for _, button := range bm.buttons {
			level := button.pin.Read()
			if level == button.pressedLevel {
				button.handler()
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(bm.cfg.PollPeriod):
		}
	}
}
