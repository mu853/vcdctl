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
		NewCmdGetVAppNetwork(),
		NewCmdGetVAppVm(),
		NewCmdGetVAppVmNetwork(),
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
					nw.FenceMode,
					nw.IsShared,
					nw.IsIpScopeInherited})
			}
			PrityPrint([]string{"Name", "Id", "Org", "Vdc", "DefaultGateway", "Dns1", "Dns2", "DnsSuffix", "FenceMode", "IsShared", "IsIpScopeInherited"}, data)
		},
	}
	return cmd
}

func NewCmdGetVAppNetwork() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vapp-network ${VAPP_NAME}",
		Short:   "Get VAppNetwork [an]",
		Aliases: []string{"an"},
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			vapps := GetVApps()
			vappNames := []string{}
			for _, vapp := range vapps {
				vappNames = append(vappNames, vapp.Name)
			}

			return vappNames, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			vapp := GetVAppByName(args[0])
			vdcNetworks := GetVdcNetwork(GetVdcId(vapp.VdcName))

			var data [][]string
			for _, nw := range GetVAppNetwork(vapp.Id) {
				IpScope := nw.Configuration.IpScopes.IpScope[0]
				var vdcNetwork OrgVdcNetwork
				for _, vdcnw := range vdcNetworks {
					if nw.Configuration.ParentNetwork.Id == vdcnw.Id {
						vdcNetwork = vdcnw
					}
				}
				data = append(data, []string{
					nw.Name,
					IpScope.IsInherited,
					IpScope.IsEnabled,
					IpScope.Gateway + "/" + IpScope.SubnetPrefixLength,
					nw.Configuration.ParentNetwork.Name,
					nw.Configuration.ParentNetwork.Id,
					vdcNetwork.FenceMode})
			}
			PrityPrint([]string{"Name", "IsInherited", "IsEnabled", "DefaultGateway", "ParentName", "ParentId", "ParentFenceMode"}, data)
		},
	}
	return cmd
}

func NewCmdGetVAppVm() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vapp-vm ${VAPP_NAME}",
		Short:   "Get VApp VMs [vm]",
		Aliases: []string{"vm"},
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			vapps := GetVApps()
			vappNames := []string{}
			for _, vapp := range vapps {
				vappNames = append(vappNames, vapp.Name)
			}

			return vappNames, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			vapp := GetVAppByName(args[0])

			var data [][]string
			for _, vm := range GetVAppVm(vapp.Id) {
				data = append(data, []string{
					vm.Name,
					vm.Urn,
					vm.Href})
			}
			PrityPrint([]string{"Name", "Urn", "Href"}, data)
		},
	}
	return cmd
}

func NewCmdGetVAppVmNetwork() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vapp-vmnetwork ${VAPP_NAME}",
		Short:   "Get VApp VM Networks [vmn]",
		Aliases: []string{"vmn"},
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			vapps := GetVApps()
			vappNames := []string{}
			for _, vapp := range vapps {
				vappNames = append(vappNames, vapp.Name)
			}

			return vappNames, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			vapp := GetVAppByName(args[0])

			var data [][]string
			for _, vm := range GetVAppVm(vapp.Id) {
				for _, nw := range vm.NetworkConnectionSection.NetworkConnection {
					data = append(data, []string{
						vm.Name,
						vm.Urn,
						strconv.Itoa(nw.NetworkConnectionIndex),
						nw.IsConnected,
						nw.NetworkAdapterType,
						nw.Name,
						nw.IpAddressAllocationMode,
						nw.IpAddress,
						nw.MACAddress})
				}
			}
			PrityPrint([]string{"Vm", "VmId", "Index", "IsConnected", "Type", "Network", "Mode", "IpAddress", "MacAddress"}, data)
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

	for i := 0; i < len(orgList.Org); i++ {
		orgList.Org[i].Id = LastOne(orgList.Org[i].Href, "/")
	}
	return orgList.Org
}

func GetOrgVdcs() []OrgVdc {
	res := client.Request("GET", "/api/query?type=adminOrgVdc", nil, nil)

	var orgVdcList OrgVdcList
	err := xml.Unmarshal(res.Body, &orgVdcList)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(orgVdcList.OrgVdc); i++ {
		orgVdcList.OrgVdc[i].Id = LastOne(orgVdcList.OrgVdc[i].Href, "/")
	}
	return orgVdcList.OrgVdc
}

func GetVApps() []VApp {
	res := client.Request("GET", "/api/admin/extension/vapps/query", nil, nil)

	var vappList VAppList
	err := xml.Unmarshal(res.Body, &vappList)
	if err != nil {
		log.Fatal(err)
	}

	var orgList []Org = GetOrgs()

	for i := 0; i < len(vappList.VApp); i++ {
		vappList.VApp[i].Id = LastOne(vappList.VApp[i].Href, "/")
		for _, org := range orgList {
			if vappList.VApp[i].OrgHref == org.Href {
				vappList.VApp[i].OrgName = org.Name
				break
			}
		}
	}
	return vappList.VApp
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
		networkType := GetVdcNetworkType(networkList.OrgVdcNetworks[i].Id)
		for _, vdc := range vdcList {
			if vdcId == vdc.Id {
				networkList.OrgVdcNetworks[i].OrgName = vdc.OrgName
				networkList.OrgVdcNetworks[i].FenceMode = networkType
				break
			}
		}
	}
	return networkList.OrgVdcNetworks
}

func GetVdcNetworkType(networkId string) string {
	res := client.Request("GET", "/api/network/"+networkId, nil, nil)

	var network Network
	err := xml.Unmarshal(res.Body, &network)
	if err != nil {
		log.Fatal(err)
	}

	return network.Configuration.FenceMode
}

func GetVdcId(vdcName string) string {
	for _, vdc := range GetOrgVdcs() {
		if vdc.Name == vdcName {
			return vdc.Id
		}
	}
	return ""
}

func GetVAppByName(vappName string) VApp {
	for _, vapp := range GetVApps() {
		if vapp.Name == vappName {
			return vapp
		}
	}
	return VApp{}
}

func GetVAppNetwork(vappId string) []Network {
	res := client.Request("GET", "/api/vApp/"+vappId+"/networkConfigSection", nil, nil)

	var networkConfigSection NetworkConfigSection
	err := xml.Unmarshal(res.Body, &networkConfigSection)
	if err != nil {
		log.Fatal(err)
	}

	return networkConfigSection.NetworkConfig
}

func GetVAppVm(vappId string) []VM {
	res := client.Request("GET", "/api/vApp/"+vappId, nil, nil)

	var vappDetails VAppDetails
	err := xml.Unmarshal(res.Body, &vappDetails)
	if err != nil {
		log.Fatal(err)
	}

	return vappDetails.VMs.VM
}
