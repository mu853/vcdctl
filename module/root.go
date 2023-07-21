package module

import (
	"github.com/spf13/cobra"
)

var (
	client         VcdClient
	config         Config
	configFilePath string
)

func GetCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vcdctl",
		Short: "vCD command-line client",
		Long:  "vCD command-line client",
	}
	cmd.AddCommand(
		NewCmdGet(),
		NewCmdPost(),
		NewCmdConfig(),
		NewCmdApi(),
	)
	cmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", defaultConfigFilePath(), "path to vcdctl config file")

	return cmd
}
