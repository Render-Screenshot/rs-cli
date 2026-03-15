package main

import (
	"os"

	"github.com/Render-Screenshot/rs-cli/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
