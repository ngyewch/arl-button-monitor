package gamepad

type Config struct {
	Buttons []ButtonConfig `json:"buttons" validate:"required,min=1,dive"`
}

type ButtonConfig struct {
	ID        string `json:"id" validate:"required"`
	GamepadID string `json:"gamepadId" validate:"required"`
	ButtonID  int    `json:"buttonId"`
}
