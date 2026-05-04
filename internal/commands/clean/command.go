package clean

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func NewKubeletCommand() *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)

	cmd := &cobra.Command{
		Use:   "kubectl clean",
		Short: "Clean Useless resource",
		Run:   cmdutil.DefaultSubCommandRun(os.Stderr),
	}

	cmd.AddCommand(NewCmdCleanRs(configFlags))
	return cmd
}
