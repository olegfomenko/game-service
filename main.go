package main

import (
	"os"

  "github.com/olegfomenko/game-service/internal/cli"

)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
