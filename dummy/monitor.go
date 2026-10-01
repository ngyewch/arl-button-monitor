package dummy

import (
	buttonmonitor "github.com/ngyewch/arl-button-monitor"
	"github.com/ngyewch/arl-button-monitor/base"
)

const (
	topicButtonPressed = "buttonPressed"
)

type ButtonMonitor struct {
	base *base.ButtonMonitor
}

func New() (*ButtonMonitor, error) {
	baseButtonMonitor, err := base.New()
	if err != nil {
		return nil, err
	}
	return &ButtonMonitor{
		base: baseButtonMonitor,
	}, nil
}

func (bm *ButtonMonitor) Close() error {
	_ = bm.base.Close()
	return nil
}

func (bm *ButtonMonitor) Subscribe() (buttonmonitor.Subscription[string], error) {
	return bm.base.Subscribe()
}
