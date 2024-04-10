package module

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "create resources",
		Args:  cobra.MaximumNArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initClient()
		},
	}
	cmd.AddCommand(
		NewCmdSetOrgVdcNetwork(),
	)
	return cmd
}

func NewCmdSetOrgVdcNetwork() *cobra.Command {
	var orgvdcName string
	var connected bool
	var edgeName string
	var distributed bool

	cmd := &cobra.Command{
		Use:     "vdc-network ${NETWORK_NAME}",
		Short:   "Connect/Disconnect VdcNetwork to/from Edge [vn]",
		Aliases: []string{"vn"},
		Args:    cobra.ExactArgs(1),
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

			if connected {
				edge, err := GetEdge(edgeName, orgvdcName)
				if err != nil {
					Fatal(err)
				}
				connectionType := "NON_DISTRIBUTED"
				if distributed {
					connectionType = "INTERNAL"
				}
				network.Connection = &ConnectionInfo{
					RouterRef:           ReferenceJson{ Urn: edge.Urn },
					ConnectionTypeValue: connectionType,
					Connected: true,
				}
				network.NetworkType = "NAT_ROUTED"
			} else {
				network.Connection = nil
				network.NetworkType = "ISOLATED"
			}

			data, err := json.Marshal(network)
			if err != nil {
				Fatal(err)
			}

			header := map[string]string{"Content-Type": "application/json"}
			client.Request("PUT", fmt.Sprintf("/cloudapi/1.0.0/orgVdcNetworks/%s", network.Urn), header, data)
		},
	}
	cmd.PersistentFlags().StringVarP(&orgvdcName, "orgvdc", "", "", "org vdc name (required)")
	cmd.PersistentFlags().BoolVarP(&connected, "connected", "", true, "connect network to edge (default true)")
	cmd.PersistentFlags().StringVarP(&edgeName, "edge", "", "", "edge name (required if connected is true)")
	cmd.PersistentFlags().BoolVarP(&distributed, "distributed", "", false, "enable distributed connection (default false)")
	cmd.MarkFlagRequired("orgvdc")

	cmd.RegisterFlagCompletionFunc("orgvdc", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		initClient()
		return GetOvdcNames(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc("edge", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		initClient()
		return GetEdgeNames(orgvdcName), cobra.ShellCompDirectiveNoFileComp
	})
	return cmd
}
