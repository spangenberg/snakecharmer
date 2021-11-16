package main

import (
	"github.com/spangenberg/snakecharmer"
	"github.com/spangenberg/snakecharmer/examples/gendocs/internal/cmd"
)

func main() {
	snakecharmer.GenDocs(cmd.NewCmdRoot())
}
