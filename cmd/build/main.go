package main

import (
	"os"

	"github.com/jimmykodes/recipes/internal/cmds"
)

func main() {
	if err := cmds.Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
