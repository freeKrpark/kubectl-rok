package images

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func NewKubeletCommand() *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)

	var allNamespaces bool

	cmd := &cobra.Command{
		Use:   "kubectl images",
		Short: "Show pods image Respoitory, tag",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(configFlags, allNamespaces)
		},
	}

	configFlags.AddFlags(cmd.PersistentFlags())
	cmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "List pods in all namespaces")

	return cmd
}
