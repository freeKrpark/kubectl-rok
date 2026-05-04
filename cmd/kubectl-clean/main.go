package main

import (
	"github.com/freekrpark/kubectl-rok/internal/cli"
	"github.com/freekrpark/kubectl-rok/internal/commands/clean"
)

func main() {
	command := clean.NewKubeletCommand()
	cli.Run(command)
}
