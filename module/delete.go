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
		NewCmdDeleteOrgVdcNetwork(),
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

func NewCmdDeleteOrgVdcNetwork() *cobra.Command {
	var orgvdcName string

	cmd := &cobra.Command{
		Use:     "vdc-network ${NETWORK_NAME}",
		Short:   "Delete VdcNetwork [vn]",
		Aliases: []string{"vn"},
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			initClient()
			networkNames := []string{}
			vdc, err := GetVdc(orgvdcName)
			if err != nil {
				Fatal(err)
			}
			for _, nw := range GetOrgVdcNetworks(vdc.Id) {
				networkNames = append(networkNames, nw.Name)
			}
			return networkNames, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			networkName := args[0]

			vdc, err := GetVdc(orgvdcName)
			if err != nil {
				Fatal(err)
			}

			network, err := GetOrgVdcNetwork(networkName, vdc.Id)
			if err != nil {
				Fatal(err)
			}
			client.Request("DELETE", fmt.Sprintf("/cloudapi/1.0.0/orgVdcNetworks/%s", network.Urn), nil, nil)
		},
	}
	cmd.PersistentFlags().StringVarP(&orgvdcName, "orgvdc", "", "", "org vdc name (required)")
	cmd.MarkFlagRequired("orgvdc")

	cmd.RegisterFlagCompletionFunc("orgvdc", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		initClient()
		return GetOvdcNames(), cobra.ShellCompDirectiveNoFileComp
	})
	return cmd
}