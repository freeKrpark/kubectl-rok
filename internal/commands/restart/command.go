package restart

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func NewKubeletCommand() *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)

	var allNamespaces bool

	cmd := &cobra.Command{
		Use:   "kubectl restart-history",
		Short: "Show pods shorted by restart count",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(configFlags, allNamespaces)
		},
	}

	configFlags.AddFlags(cmd.PersistentFlags())
	cmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "List pods in all namespaces")
	return cmd
}
