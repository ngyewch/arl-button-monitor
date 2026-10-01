package main

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/jwalton/go-supportscolor"
	"github.com/phsym/console-slog"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/host/v3"
)

func init() {
	initPeriph()
	initLogging()
}

func initLogging() {
	logLevel := slog.LevelInfo
	_ = logLevel.UnmarshalText([]byte(os.Getenv("LOG_LEVEL")))
	addSource := os.Getenv("LOG_ADD_SOURCE") == "true"
	noColor := os.Getenv("NO_COLOR") == "true"

	supportsColor := supportscolor.Stderr().SupportsColor
	consoleHandler := console.NewHandler(os.Stderr, &console.HandlerOptions{
		Level:     logLevel,
		AddSource: addSource,
		NoColor:   noColor || !supportsColor,
	})
	logger := slog.New(consoleHandler)
	slog.SetDefault(logger)
}

func initPeriph() {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	err := doInitPeriph()
	if err != nil {
		log.SetOutput(os.Stderr)
		slog.Error("init error",
			slog.Any("err", err),
		)
		os.Exit(1)
	}
}

func doInitPeriph() error {
	_, err := host.Init()
	if err != nil {
		return err
	}

	_, err = driverreg.Init()
	if err != nil {
		return err
	}

	return nil
}
