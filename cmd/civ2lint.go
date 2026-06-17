package cmd

import (
	"flag"
	"os"

	"go.uber.org/zap"

	"github.com/leonsp/civ2lint/lib"
)

func Init() {
	var c lib.Config
	var usage bool

	flag.StringVar(&c.Path, "path", ".", "Path to the game, mod, or scenario directory")
	flag.BoolVar(&usage, "help", false, "Print usage instructions")
	flag.Parse()

	if usage {
		flag.Usage()
		return
	}
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }() // flushes buffer, if any
	sugar := logger.Sugar()

	sugar.Info("Logger initialized")

	cl := lib.New(c, sugar)
	err := cl.Lint()
	if err != nil {
		sugar.Error(err)
		os.Exit(1)
	}
}
