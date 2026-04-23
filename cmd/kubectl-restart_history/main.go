package main

import (
	"github.com/freekrpark/kubectl-rok/internal/cli"
	"github.com/freekrpark/kubectl-rok/internal/commands/restart"
)

func main() {
	command := restart.NewKubeletCommand()
	cli.Run(command)
}
