package shell

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func NewKubeletCommand() *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)

	var container string
	var debug bool
	var tool string
	cmd := &cobra.Command{
		Use:   "kubectl-shell [pod-name]",
		Short: "Drop into a shell inside a pod",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(configFlags, args[0], container, debug, tool)
		},
	}

	configFlags.AddFlags(cmd.PersistentFlags())
	cmd.Flags().StringVarP(&container, "container", "c", "", "Container name (default: first non-sidecar)")
	cmd.Flags().BoolVar(&debug, "debug", false, "Use ephemeral debug container (for distroless images)")
	cmd.Flags().StringVar(&tool, "tool", "busybox", "Debug image: busybox|netshoot|ubuntu|alpine")
	return cmd
}
