package clean

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func NewCmdCleanRs(configFlags *genericclioptions.ConfigFlags) *cobra.Command {

	var allNamespaces bool
	cmd := &cobra.Command{
		Use:   "rs",
		Short: "Delete ReplicaSets with 0 desired replicas",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(configFlags, allNamespaces)
		},
	}

	configFlags.AddFlags(cmd.PersistentFlags())
	cmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "Clean across all namespaces")

	return cmd
}
