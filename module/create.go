package module

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
)

func NewCmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create resources",
		Args:  cobra.MaximumNArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initClient()
		},
	}
	cmd.AddCommand(
		NewCmdCreateOrg(),
		NewCmdCreateOrgVdc(),
		NewCmdCreateOrgVdcNetwork(),
		NewCmdCreateVApp(),
		NewCmdCreateVAppNetwork(),
	)
	return cmd
}

func NewCmdCreateOrg() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org",
		Short: "Create Organization",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			orgName := args[0]

			org := struct {
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
			}{ Name: orgName, DisplayName: orgName }

			data, err := json.Marshal(org)
			if err != nil {
				Fatal(nil)
			}

			header := map[string]string{ "Content-Type": "application/json" }
			client.Request("POST", "/cloudapi/1.0.0/orgs", header, data)
		},
	}
	return cmd
}

func NewCmdCreateOrgVdc() *cobra.Command {
	var orgName string
	var providerVdcName string
	var storagePolicyName string
	var networkPoolName string

	cmd := &cobra.Command{
		Use:     "orgvdc",
		Aliases: []string{"vdc"},
		Short:   "Create Org VDC [vdc]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			vdcName := args[0]

			type CapacityWithUsageType struct {
				Units           string `xml:"Units"`
				Allocated       int    `xml:"Allocated"`
				Limit           int    `xml:"Limit"`
				Reserved        int    `xml:"Reserved"`
				Used            int    `xml:"Used"`
				ReservationUsed int    `xml:"ReservationUsed"`
			}
			type ComputeCapacity struct {
				Cpu    CapacityWithUsageType `xml:"Cpu"`
				Memory CapacityWithUsageType `xml:"Memory"`
			}
			type VdcStorageProfile struct {
				Enabled                   bool         `xml:"Enabled"`
				Units                     string       `xml:"Units"`
				Limit                     int          `xml:"Limit"`
				Default                   bool         `xml:"Default"`
				ProviderVdcStorageProfile VdcReference `xml:"ProviderVdcStorageProfile"`
			}
			type CreateVdcParams struct {
				Xmlns                    string            `xml:"xmlns,attr"`
				XmlnsExtension           string            `xml:"xmlns:extension_v1.5,attr"`
				Name                     string            `xml:"name,attr"`
				AllocationModel          string            `xml:"AllocationModel"`
				ComputeCapacity          ComputeCapacity   `xml:"ComputeCapacity"`
				IsEnabled                bool              `xml:"IsEnabled"`
				VdcStorageProfile        VdcStorageProfile `xml:"VdcStorageProfile"`
				ResourceGuaranteedMemory float64           `xml:"ResourceGuaranteedMemory"`
				ResourceGuaranteedCpu    float64           `xml:"ResourceGuaranteedCpu"`
				IsThinProvision          bool              `xml:"IsThinProvision"`
				NetworkPoolReference     VdcReference      `xml:"NetworkPoolReference"`
				ProviderVdcReference     VdcReference      `xml:"ProviderVdcReference"`
				UsesFastProvisioning     bool              `xml:"UsesFastProvisioning"`
				VmDiscoveryEnabled       bool              `xml:"VmDiscoveryEnabled"`
				IncludeMemoryOverhead    bool              `xml:"IncludeMemoryOverhead"`
			}

			newVdc := CreateVdcParams{
				Xmlns: "http://www.vmware.com/vcloud/v1.5",
				XmlnsExtension: "http://www.vmware.com/vcloud/extension/v1.5",
				Name: vdcName,
				AllocationModel: "Flex",
				ComputeCapacity: ComputeCapacity{
					Cpu: CapacityWithUsageType{
						Limit: 0,
						Reserved: 0,
						Units: "MHz",
					},
					Memory: CapacityWithUsageType{
						Limit: 0,
						Reserved: 0,
						Units: "MB",
					},
				},
				IsEnabled: true,
				VdcStorageProfile: VdcStorageProfile{
					Enabled: true,
					Units: "MB",
					Default: true,
					ProviderVdcStorageProfile: GetStorageProfile(storagePolicyName, providerVdcName),
				},
				ResourceGuaranteedMemory: 0.0,
				ResourceGuaranteedCpu: 0.0,
				IsThinProvision: true,
				NetworkPoolReference: GetNetworkPool(networkPoolName),
				ProviderVdcReference: GetProviderVdc(providerVdcName),
				UsesFastProvisioning: false,
				VmDiscoveryEnabled: false,
				IncludeMemoryOverhead: false,
			}
			org := GetOrg(orgName)

			data, err := xml.Marshal(newVdc)
			if err != nil {
				Fatal(err)
			}

			header := map[string]string{ "Content-Type": "application/vnd.vmware.admin.createVdcParams+xml" }
			res := client.Request("POST", fmt.Sprintf("/api/admin/org/%s/vdcsparams", org.Id), header, data)
			fmt.Println(string(res.Body))
		},
	}
	cmd.PersistentFlags().StringVarP(&orgName, "org", "o", "", "org name")
	cmd.PersistentFlags().StringVarP(&providerVdcName, "provider-vdc", "p", "", "provider vdc name")
	cmd.PersistentFlags().StringVarP(&storagePolicyName, "storage-policy", "s", "", "storage policy name")
	cmd.PersistentFlags().StringVarP(&networkPoolName, "network-pool", "n", "", "network pool name")
	return cmd
}

func NewCmdCreateOrgVdcNetwork() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vdc-network ${VDC_NAME}",
		Short:   "Create VdcNetwork [vn]",
		Aliases: []string{"vn"},
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			initClient()
			orgVdcs := GetOrgVdcs()
			vdcNames := []string{}
			for _, vdc := range orgVdcs {
				vdcNames = append(vdcNames, vdc.Name)
			}

			return vdcNames, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			vcdName := args[0]
			orgVdc, err := GetVdc(vcdName)
			if err != nil {
				Fatal(err)
			}
			var data [][]string
			for _, nw := range GetOrgVdcNetwork(orgVdc.Id) {
				ipScope := nw.Configuration.IpScopes.IpScope[0]
				data = append(data, []string{
					nw.Name,
					nw.Id,
					orgVdc.OrgName,
					vcdName,
					ipScope.Gateway + "/" + ipScope.SubnetPrefixLength,
					ipScope.Dns1,
					ipScope.Dns2,
					ipScope.DnsSuffix,
					nw.Configuration.FenceMode,
					nw.IsShared,
					ipScope.IsInherited})
			}
			PrityPrint([]string{"Name", "Id", "Org", "Vdc", "DefaultGateway", "Dns1", "Dns2", "DnsSuffix", "FenceMode", "IsShared", "IsIpScopeInherited"}, data)
		},
	}
	return cmd
}

func NewCmdCreateVApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vapp",
		Aliases: []string{"a"},
		Short:   "Create VApp [a]",
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

func NewCmdCreateVAppNetwork() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vapp-network ${VAPP_NAME}",
		Short:   "Create VAppNetwork [an]",
		Aliases: []string{"an"},
		Args:    cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			initClient()
			vapps := GetVApps()
			vappNames := []string{}
			for _, vapp := range vapps {
				vappNames = append(vappNames, vapp.Name)
			}

			return vappNames, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			vapp, err := GetVAppByName(args[0])
			if err != nil {
				Fatal(err)
			}
			orgVdc, err := GetVdc(vapp.VdcName)
			if err != nil {
				Fatal(err)
			}
			vdcNetworks := GetOrgVdcNetwork(orgVdc.Id)

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
					vdcNetwork.Configuration.FenceMode})
			}
			PrityPrint([]string{"Name", "IsInherited", "IsEnabled", "DefaultGateway", "ParentName", "ParentId", "ParentFenceMode"}, data)
		},
	}
	return cmd
}

func GetOrgs_() []Org {
	res := client.Request("GET", "/api/org", nil, nil)
	var orgList OrgList
	err := xml.Unmarshal(res.Body, &orgList)
	if err != nil {
		Fatal(err)
	}

	for i := 0; i < len(orgList.Org); i++ {
		orgList.Org[i].Id = LastOne(orgList.Org[i].Href, "/")
	}

	orgs := orgList.Org
	sort.Slice(orgs, func(i, j int) bool {
		return orgs[i].Name < orgs[j].Name
	})
	return orgs
}
