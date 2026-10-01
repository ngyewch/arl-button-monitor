package ft232h

import (
	"fmt"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3/ftdi"
)

func findFT232H() *ftdi.FT232H {
	all := ftdi.All()
	for _, d := range all {
		ft232h, ok := d.(*ftdi.FT232H)
		if ok {
			return ft232h
		}
	}
	return nil
}

func pullFromString(s string) (gpio.Pull, error) {
	switch s {
	case "PullNoChange":
		return gpio.PullNoChange, nil
	case "Float":
		return gpio.Float, nil
	case "PullDown":
		return gpio.PullDown, nil
	case "PullUp":
		return gpio.PullUp, nil
	}
	return gpio.PullNoChange, fmt.Errorf("invalid pull value: %s", s)
}

func edgeFromString(s string) (gpio.Edge, error) {
	switch s {
	case "NoEdge":
		return gpio.NoEdge, nil
	case "RisingEdge":
		return gpio.RisingEdge, nil
	case "FallingEdge":
		return gpio.FallingEdge, nil
	case "BothEdges":
		return gpio.BothEdges, nil
	}
	return gpio.NoEdge, fmt.Errorf("invalid edge value: %s", s)
}

func levelFromString(s string) (gpio.Level, error) {
	switch s {
	case "Low":
		return gpio.Low, nil
	case "High":
		return gpio.High, nil
	}
	return gpio.Low, fmt.Errorf("invalid level value: %s", s)
}
