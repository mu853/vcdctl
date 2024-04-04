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
	cmd.AddCommand(
		NewCmdDeleteOrg(),
	)
	return cmd
}

func NewCmdDeleteOrg() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org",
		Short: "Delete Organization",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			org := GetOrg(args[0])
			client.Request("DELETE", fmt.Sprintf("/cloudapi/1.0.0/orgs/urn:vcloud:org:%s", org.Id), nil, nil)
		},
	}
	return cmd
}
