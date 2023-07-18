package module

import (
	"encoding/xml"
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get resources or exec get api",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				return
			}
			api := args[0]
			res := client.Request("GET", api, nil, nil)
			fmt.Println(string(res.Body))
		},
	}
	cmd.AddCommand(
		NewCmdGetOrg(),
		NewCmdGetOrgVdc(),
		NewCmdGetVApp(),
		NewCmdGetVdcNetwork(),
	)
	return cmd
}

func NewCmdGetOrg() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org",
		Short: "Get Organization",
		Run: func(cmd *cobra.Command, args []string) {
			header := []string{"Name", "Id", "href"}
			var data [][]string
			for _, org := range GetOrgs() {
				data = append(data, []string{org.Name, org.Id, org.Href})
			}
			PrityPrint(header, data)
		},
	}
	return cmd
}

func NewCmdGetOrgVdc() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "orgvdc",
		Aliases: []string{"vdc"},
		Short:   "Get Org VDC [vdc]",
		Run: func(cmd *cobra.Command, args []string) {
			header := []string{"Name", "Id", "IsEnabled", "Org", "ProviderVdc", "Vc", "NetworkType", "VApps", "VMs", "VAppTemplates"}
			var data [][]string
			for _, vdc := range GetOrgVdcs() {
				data = append(data, []string{
					vdc.Name,
					vdc.Id,
					vdc.IsEnabled,
					vdc.OrgName,
					vdc.ProviderVdcName,
					vdc.VcName,
					vdc.NetworkProviderType,
					strconv.Itoa(vdc.NumberOfVApps),
					strconv.Itoa(vdc.NumberOfVMs),
					strconv.Itoa(vdc.NumberOfVAppTemplates)})
			}
			PrityPrint(header, data)
		},
	}
	return cmd
}

func NewCmdGetVApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vapp",
		Short: "Get VApp",
		Run: func(cmd *cobra.Command, args []string) {
			var data [][]string
			for _, vapp := range GetVApps() {
				data = append(data, []string{
					vapp.Name,
					vapp.Id,
					vapp.IsEnabled,
					vapp.Status,
					vapp.OrgName,
					vapp.VdcName,
					strconv.Itoa(vapp.NumberOfVMs),
					vapp.TaskStatusName,
					vapp.TaskStatus})
			}
			PrityPrint([]string{"Name", "Id", "IsEnabled", "Status", "Org", "Vdc", "VMs", "TaskStatusName", "TaskStatus"}, data)
		},
	}
	return cmd
}

func NewCmdGetVdcNetwork() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vdc-network ${VDC_NAME}",
		Short:   "Get VdcNetwork [vn]",
		Aliases: []string{"vn"},
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			orgVdcs := GetOrgVdcs()
			vdcNames := []string{}
			for _, vdc := range orgVdcs {
				vdcNames = append(vdcNames, vdc.Name)
			}

			return vdcNames, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			var data [][]string
			for _, nw := range GetVdcNetwork(GetVdcId(args[0])) {
				data = append(data, []string{
					nw.Name,
					nw.Id,
					nw.OrgName,
					nw.VdcName,
					nw.DefaultGateway + "/" + nw.SubnetPrefixLength,
					nw.Dns1,
					nw.Dns2,
					nw.DnsSuffix,
					nw.IsShared,
					nw.IsIpScopeInherited})
			}
			PrityPrint([]string{"Name", "Id", "Org", "Vdc", "DefaultGateway", "Dns1", "Dns2", "DnsSuffix", "IsShared", "IsIpScopeInherited"}, data)
		},
	}
	return cmd
}

func GetOrgs() []Org {
	res := client.Request("GET", "/api/org", nil, nil)
	var orgList OrgList
	err := xml.Unmarshal(res.Body, &orgList)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(orgList.Orgs); i++ {
		orgList.Orgs[i].Id = LastOne(orgList.Orgs[i].Href, "/")
	}
	return orgList.Orgs
}

func GetOrgVdcs() []OrgVdc {
	res := client.Request("GET", "/api/query?type=adminOrgVdc", nil, nil)

	var orgVdcList OrgVdcList
	err := xml.Unmarshal(res.Body, &orgVdcList)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(orgVdcList.OrgVdcs); i++ {
		orgVdcList.OrgVdcs[i].Id = LastOne(orgVdcList.OrgVdcs[i].Href, "/")
	}
	return orgVdcList.OrgVdcs
}

func GetVApps() []VApp {
	res := client.Request("GET", "/api/admin/extension/vapps/query", nil, nil)

	var vappList VAppList
	err := xml.Unmarshal(res.Body, &vappList)
	if err != nil {
		log.Fatal(err)
	}

	var orgList []Org = GetOrgs()

	for i := 0; i < len(vappList.VApps); i++ {
		vappList.VApps[i].Id = LastOne(vappList.VApps[i].Href, "/")
		for _, org := range orgList {
			if vappList.VApps[i].OrgHref == org.Href {
				vappList.VApps[i].OrgName = org.Name
				break
			}
		}
	}
	return vappList.VApps
}

func GetVdcNetwork(vdcId string) []OrgVdcNetwork {
	res := client.Request("GET", "/api/admin/vdc/"+vdcId+"/networks", nil, nil)

	var networkList OrgVdcNetworkList
	err := xml.Unmarshal(res.Body, &networkList)
	if err != nil {
		log.Fatal(err)
	}

	var vdcList []OrgVdc = GetOrgVdcs()

	for i := 0; i < len(networkList.OrgVdcNetworks); i++ {
		networkList.OrgVdcNetworks[i].Id = LastOne(networkList.OrgVdcNetworks[i].Href, "/")
		vdcId = LastOne(networkList.OrgVdcNetworks[i].VdcHref, "/")
		for _, vdc := range vdcList {
			if vdcId == vdc.Id {
				networkList.OrgVdcNetworks[i].OrgName = vdc.OrgName
				break
			}
		}
	}
	return networkList.OrgVdcNetworks
}

func GetVdcId(vdcName string) string {
	for _, vdc := range GetOrgVdcs() {
		if vdc.Name == vdcName {
			return vdc.Id
		}
	}
	return ""
}
