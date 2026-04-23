package main

import (
	"github.com/freekrpark/kubectl-rok/internal/cli"
	"github.com/freekrpark/kubectl-rok/internal/commands/drift"
)

func main() {
	command := drift.NewKubeletCommand()
	cli.Run(command)
}
