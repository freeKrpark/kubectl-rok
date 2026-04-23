package main

import (
	"github.com/freekrpark/kubectl-rok/internal/cli"
	"github.com/freekrpark/kubectl-rok/internal/commands/images"
)

func main() {
	command := images.NewKubeletCommand()
	cli.Run(command)
}
