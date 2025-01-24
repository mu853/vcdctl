package module

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewCmdPut() *cobra.Command {
	var fileName string
	var header string
	cmd := &cobra.Command{
		Use:   "put",
		Short: "exec put api",
		Args:  cobra.ExactArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initClient()
		},
		Run: func(cmd *cobra.Command, args []string) {
			api := args[0]

			var data []byte
			if fileName != "" {
				data = ReadRequestData(fileName)
			}

			var header_map map[string]string
			if header != "" {
				header_map = map[string]string{}
				entry := strings.Split(header, ": ")
				header_map[entry[0]] = entry[1]
			}
			res := client.Request("PUT", api, header_map, data)
			fmt.Println(string(res.Body))
		},
	}
	cmd.Flags().StringVarP(&fileName, "filename", "f", "", "file name for send data(xml)")
	cmd.Flags().StringVarP(&header, "header", "", "", "additional header (cf. \"Content-Type: application/vnd.vmware.vcloud.vm+xml\"")

	cmd.AddCommand(
		NewCmdPutVAppLease(),
	)
	return cmd
}

func NewCmdPutVAppLease() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lease",
		Short: "Extend vApp Deployment Lease",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			initClient()
			return GetVAppNames(), cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}
			vapp, err := GetVAppByName(args[0])
			if err != nil {
				Fatal(err)
			}
			vappLease := GetVAppLease(vapp.Id)
			vappLease.DeploymentLeaseExpiration = ""
			data, err := xml.Marshal(vappLease)
			Log(string(data))
			if err != nil {
				Fatal(err)
			}
			header := map[string]string{"Content-Type": "application/vnd.vmware.vcloud.leaseSettingsSection+xml"}
			client.Request("PUT", fmt.Sprintf("/api/vApp/%s/leaseSettingsSection", vapp.Id), header, data)
		},
	}
	return cmd
}
