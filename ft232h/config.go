package ft232h

import "time"

type Config struct {
	PollPeriod     time.Duration  `json:"pollPeriod" validate:"required"`
	DebouncePeriod time.Duration  `json:"debouncePeriod" validate:"required"`
	Buttons        []ButtonConfig `json:"buttons" validate:"required,min=1,dive"`
}

type ButtonConfig struct {
	ID           string `json:"id" validate:"required"`
	GPIOPinName  string `json:"gpioPinName" validate:"required"`
	PressedLevel string `json:"pressedLevel" validate:"required,oneof=Low High"`
}
