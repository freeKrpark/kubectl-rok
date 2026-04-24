package main

import (
	"github.com/freekrpark/kubectl-rok/internal/cli"
	"github.com/freekrpark/kubectl-rok/internal/commands/shell"
)

func main() {
	command := shell.NewKubeletCommand()
	cli.Run(command)
}
