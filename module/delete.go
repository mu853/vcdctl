package module

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "exec delete api",
		Args:  cobra.ExactArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initClient()
		},
		Run: func(cmd *cobra.Command, args []string) {
			api := args[0]
			res := client.Request("DELETE", api, nil, nil)
			fmt.Println(string(res.Body))
		},
	}
	return cmd
}
