package main

import (
	"flag"
	"os"

	"github.com/x0k/veterinary-clinic-backend/internal/app"
	"github.com/x0k/veterinary-clinic-backend/internal/config"
)

var (
	config_path string
)

func init() {
	flag.StringVar(&config_path, "config", os.Getenv("CONFIG_PATH"), "Config path")
	flag.Parse()
}

func main() {
	app.Run(config.MustLoad(config_path))
}
