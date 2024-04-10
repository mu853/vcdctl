package module

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

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
		NewCmdCreateEdge(),
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
			}{Name: orgName, DisplayName: orgName}

			data, err := json.Marshal(org)
			if err != nil {
				Fatal(nil)
			}

			header := map[string]string{"Content-Type": "application/json"}
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

			if orgName == "" {
				orgName = client.site.OrgName
			}
			if orgName == "" {
				Fatal("org name not specified")
			}
			org := GetOrg(orgName)

			storageProfile, err := GetStorageProfile(storagePolicyName, providerVdcName)
			if err != nil {
				Fatal(err)
			}
			networkPool, err := GetNetworkPool(networkPoolName)
			if err != nil {
				Fatal(err)
			}
			providerVdc, err := GetProviderVdc(providerVdcName)
			if err != nil {
				Fatal(err)
			}

			newVdc := struct {
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
				NetworkPoolReference     Reference         `xml:"NetworkPoolReference"`
				ProviderVdcReference     Reference         `xml:"ProviderVdcReference"`
				UsesFastProvisioning     bool              `xml:"UsesFastProvisioning"`
				VmDiscoveryEnabled       bool              `xml:"VmDiscoveryEnabled"`
				IncludeMemoryOverhead    bool              `xml:"IncludeMemoryOverhead"`
			}{
				Xmlns:           "http://www.vmware.com/vcloud/v1.5",
				XmlnsExtension:  "http://www.vmware.com/vcloud/extension/v1.5",
				Name:            vdcName,
				AllocationModel: "Flex",
				ComputeCapacity: ComputeCapacity{
					Cpu: CapacityWithUsageType{
						Limit:    0,
						Reserved: 0,
						Units:    "MHz",
					},
					Memory: CapacityWithUsageType{
						Limit:    0,
						Reserved: 0,
						Units:    "MB",
					},
				},
				IsEnabled: true,
				VdcStorageProfile: VdcStorageProfile{
					Enabled:                   true,
					Units:                     "MB",
					Default:                   true,
					ProviderVdcStorageProfile: storageProfile,
				},
				ResourceGuaranteedMemory: 0.0,
				ResourceGuaranteedCpu:    0.0,
				IsThinProvision:          true,
				NetworkPoolReference:     networkPool,
				ProviderVdcReference:     providerVdc,
				UsesFastProvisioning:     false,
				VmDiscoveryEnabled:       false,
				IncludeMemoryOverhead:    false,
			}

			data, err := xml.Marshal(newVdc)
			if err != nil {
				Fatal(err)
			}

			header := map[string]string{"Content-Type": "application/vnd.vmware.admin.createVdcParams+xml"}
			res := client.Request("POST", fmt.Sprintf("/api/admin/org/%s/vdcsparams", org.Id), header, data)
			fmt.Println(string(res.Body))
		},
	}
	cmd.PersistentFlags().StringVarP(&orgName, "org", "", "", "org name")
	cmd.PersistentFlags().StringVarP(&providerVdcName, "provider-vdc", "", "", "provider vdc name (required)")
	cmd.PersistentFlags().StringVarP(&storagePolicyName, "storage-policy", "", "", "storage policy name (required)")
	cmd.PersistentFlags().StringVarP(&networkPoolName, "network-pool", "", "", "network pool name (required)")
	cmd.MarkFlagRequired("provider-vdc")
	cmd.MarkFlagRequired("storage-policy")
	cmd.MarkFlagRequired("network-pool")
	return cmd
}

func NewCmdCreateOrgVdcNetwork() *cobra.Command {
	var orgvdcName string
	var networkType string
	var gatewayName string
	var gatewayCidr string
	var distributed bool
	var externalNetworkName string

	cmd := &cobra.Command{
		Use:     "vdc-network ${NETWORK_NAME}",
		Short:   "Create VdcNetwork [vn]",
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

			//since omitempty only works for predefined types, using pointer type
			newVdcNetwork := OrgVdcNetworkJson{
				Name:        networkName,
				NetworkType: networkType,
				OwnerRef:    ReferenceJson{Urn: fmt.Sprintf("urn:vcloud:vdc:%s", vdc.Id)},
			}

			if networkType == "DIRECT" {
				externalNetwork, err := GetExternalNetwork(externalNetworkName)
				if err != nil {
					Fatal(err)
				}
				newVdcNetwork.Shared = true
				newVdcNetwork.ParentNetworkId = &ReferenceJson{Urn: externalNetwork.Urn}
			} else {
				gatewayCidrArr := strings.Split(gatewayCidr, "/")
				if len(gatewayCidrArr) != 2 {
					Fatal(fmt.Sprintf("gateway cidr [%s] is invalid", gatewayCidr))
				}
				prefixLen, err := strconv.Atoi(gatewayCidrArr[1])
				if err != nil {
					Fatal(err)
				}

				newVdcNetwork.Subnets = &Subnets{
					Values: []Subnet{
						{
							Gateway:      gatewayCidrArr[0],
							PrefixLength: prefixLen,
							DnsSuffix:    "",
							DnsServer1:   "",
							DnsServer2:   "",
						},
					},
				}
			}

			if networkType == "NAT_ROUTED" {
				edge, err := GetEdge(gatewayName, orgvdcName)
				if err != nil {
					Fatal(err)
				}
				connectionType := "NON_DISTRIBUTED"
				if distributed {
					connectionType = "INTERNAL"
				}
				newVdcNetwork.Connection = &ConnectionInfo{
					RouterRef:           ReferenceJson{ Urn: edge.Urn },
					ConnectionTypeValue: connectionType,
					Connected: true,
				}
			}

			data, err := json.Marshal(newVdcNetwork)
			if err != nil {
				Fatal(err)
			}

			header := map[string]string{"Content-Type": "application/json"}
			res := client.Request("POST", "/cloudapi/1.0.0/orgVdcNetworks", header, data)
			fmt.Println(string(res.Body))
		},
	}
	cmd.PersistentFlags().StringVarP(&orgvdcName, "orgvdc", "", "", "org vdc name (required)")
	cmd.PersistentFlags().StringVarP(&networkType, "type", "", "", "network type (NAT_ROUTED | ISOLATED | DIRECT) (required)")
	cmd.PersistentFlags().StringVarP(&gatewayCidr, "cidr", "", "", "gateway cidr")
	cmd.PersistentFlags().StringVarP(&gatewayName, "gateway", "", "", "gateway name (NAT_ROUTED only)")
	cmd.PersistentFlags().BoolVarP(&distributed, "distributed", "", false, "enable distributed connection (NAT_ROUTED only, default false)")
	cmd.PersistentFlags().StringVarP(&externalNetworkName, "external-network", "", "", "external network name (DIRECT only)")
	cmd.MarkFlagRequired("orgvdc")
	cmd.MarkFlagRequired("type")

	cmd.RegisterFlagCompletionFunc("orgvdc", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		initClient()
		return GetOvdcNames(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"NAT_ROUTED", "ISOLATED", "DIRECT"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc("gateway", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		initClient()
		return GetEdgeNames(orgvdcName), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc("external-network", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		initClient()
		externalNetworkNames := []string{}
		for _, nw := range GetExternalNetworks() {
			externalNetworkNames = append(externalNetworkNames, nw.Name)
		}
		return externalNetworkNames, cobra.ShellCompDirectiveNoFileComp
	})
	return cmd
}

func NewCmdCreateEdge() *cobra.Command {
	var orgvdcName string
	var providerGatewayName string
	//var primaryIp string
	//var ipRange string
	cmd := &cobra.Command{
		Use:     "edge",
		Aliases: []string{"e"},
		Short:   "Create Edge [e]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			edgeName := args[0]
			if _, err := GetEdge(edgeName, orgvdcName); err == nil {
				Fatal(fmt.Sprintf("%s is already exist", edgeName))
			}

			vdc, err := GetVdc(orgvdcName)
			if err != nil {
				Fatal(err)
			}

			providerGateway, err := GetProviderGateway(providerGatewayName)
			if err != nil {
				Fatal(err)
			}
			providerGatewaySubnet := providerGateway.Subnets.Values[0]
			routerLink := EdgeGatewayUplink{
				UplinkId:   providerGateway.Urn,
				UplinkName: providerGateway.Name,
				Subnets: EdgeGatewayUplinkSubnets{
					Values: []EdgeGatewayUplinkSubnet{
						{
							GatewayAddress:       providerGatewaySubnet.GatewayAddress,
							PrefixLength:         providerGatewaySubnet.PrefixLength,
							DnsSuffix:            providerGatewaySubnet.DnsServer1,
							DnsServer1:           providerGatewaySubnet.DnsServer2,
							AutoAllocateIpRanges: true,
							TotalIpCount:         1,
						},
					},
				},
				Dedicated:     false,
				Connected:     false,
				UsingIpSpace:  false,
				VrfLiteBacked: false,
				BackingType:   "NSXT_TIER0",
			}

			newEdge := struct {
				Name                                     string              `json:"name"`
				EdgeGatewayUplinks                       []EdgeGatewayUplink `json:"edgeGatewayUplinks"`
				DistributedRoutingEnabled                bool                `json:"distributedRoutingEnabled"`
				NonDistributedRoutingEnabled             bool                `json:"nonDistributedRoutingEnabled"`
				ServiceNetworkDefinition                 string              `json:"serviceNetworkDefinition"`
				DistributedRouterUplinkNetworkDefinition string              `json:"distributedRouterUplinkNetworkDefinition"`
				OwnerRef                                 ReferenceJson       `json:"ownerRef"`
			}{
				Name:                                     edgeName,
				EdgeGatewayUplinks:                       []EdgeGatewayUplink{routerLink},
				DistributedRoutingEnabled:                false,
				NonDistributedRoutingEnabled:             true,
				ServiceNetworkDefinition:                 "192.168.255.225/27",
				DistributedRouterUplinkNetworkDefinition: "",
				OwnerRef:                                 ReferenceJson{Urn: fmt.Sprintf("urn:vcloud:vdc:%s", vdc.Id)},
			}

			data, err := json.Marshal(newEdge)
			if err != nil {
				Fatal(err)
			}

			header := map[string]string{"Content-Type": "application/json"}
			res := client.Request("POST", "/cloudapi/1.0.0/edgeGateways", header, data)
			fmt.Println(string(res.Body))
		},
	}
	cmd.PersistentFlags().StringVarP(&orgvdcName, "orgvdc", "", "", "org vdc name (required)")
	cmd.PersistentFlags().StringVarP(&providerGatewayName, "provider-gateway", "", "", "provider gateway name (required)")
	//cmd.PersistentFlags().StringVarP(&primaryIp, "primary-ip", "", "", "primary ip address")
	//cmd.PersistentFlags().StringVarP(&ipRange, "ip-range", "", "", "ip range")
	cmd.MarkFlagRequired("orgvdc")
	cmd.MarkFlagRequired("provider-gateway")

	cmd.RegisterFlagCompletionFunc("orgvdc", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		initClient()
		return GetOvdcNames(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc("provider-gateway", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		initClient()
		providerGatewayNames := []string{}
		for _, gw := range GetProviderGateways() {
			providerGatewayNames = append(providerGatewayNames, gw.Name)
		}
		return providerGatewayNames, cobra.ShellCompDirectiveNoFileComp
	})
	return cmd
}

func NewCmdCreateVApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vapp",
		Aliases: []string{"a"},
		Short:   "Create VApp [a]",
		Run: func(cmd *cobra.Command, args []string) {
			//
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
			//
		},
	}
	return cmd
}
