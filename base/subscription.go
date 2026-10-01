package base

import "github.com/F2077/go-pubsub/pubsub"

type customSubscription[T any] struct {
	subscriber   *pubsub.Subscriber[T]
	subscription *pubsub.Subscription[T]
}

func (subscription *customSubscription[T]) Close() error {
	_ = subscription.subscription.Close()
	_ = subscription.subscriber.Close()
	return nil
}

func (subscription *customSubscription[T]) Ch() <-chan T {
	return subscription.subscription.Ch
}

func (subscription *customSubscription[T]) ErrCh() <-chan error {
	return subscription.subscription.ErrCh
}
