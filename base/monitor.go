package base

import (
	"github.com/F2077/go-pubsub/pubsub"
	buttonmonitor "github.com/ngyewch/arl-button-monitor"
)

const (
	topicButtonPressed = "buttonPressed"
)

type ButtonMonitor struct {
	broker    *pubsub.Broker[string]
	publisher *pubsub.Publisher[string]
}

func New() (*ButtonMonitor, error) {
	broker, err := pubsub.NewBroker[string]()
	if err != nil {
		return nil, err
	}
	publisher := pubsub.NewPublisher[string](broker)
	bm := &ButtonMonitor{
		broker:    broker,
		publisher: publisher,
	}
	return bm, nil
}

func (bm *ButtonMonitor) Close() error {
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

func (bm *ButtonMonitor) Publish(buttonId string) error {
	return bm.publisher.Publish(topicButtonPressed, buttonId)
}
