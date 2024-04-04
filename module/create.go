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
	var orgvdcName string
	var networkType string
	var gatewayName string
	var gatewayCidr string

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

			type Reference struct {
				Id string `json:"id"`
			}

			type ConnectionInfo struct {
				RouterRef           Reference `json:"routerRef"`
				ConnectionTypeValue string    `json:"connectionTypeValue"` //INTERNAL
			}

			type Subnet struct {
				Gateway         string `json:"gateway"`
				PrefixLength    int    `json:"prefixLength"`
				DnsSuffix       string `json:"dnsSuffix"`
				DnsServer1      string `json:"dnsServer1"`
				DnsServer2      string `json:"dnsServer2"`
			}

			type Subnets struct {
				Values []Subnet `json:"values"`
			}

			//since omitempty doesn't work for predefined types, using pointer type
			type CreateVdcNetworkParams struct {
				Name            string          `json:"name"`
				Subnets         *Subnets        `json:"subnets,omitempty"`
				NetworkType     string          `json:"networkType"`
				Connection      *ConnectionInfo `json:"connection,omitempty"`
				OwnerRef        Reference       `json:"ownerRef"`
				Shared          bool            `json:"shared,omitempty"`
				ParentNetworkId *Reference      `json:"parentNetworkId,omitempty"`
			}

			vdc, err := GetVdc(orgvdcName)
			if err != nil {
				Fatal(err)
			}

			newVdcNetwork := CreateVdcNetworkParams{
				Name: networkName,
				NetworkType: networkType,
				OwnerRef: Reference{ Id: fmt.Sprintf("urn:vcloud:vdc:%s", vdc.Id) },
			}

			if networkType == "DIRECT" {
				newVdcNetwork.Shared = true
				newVdcNetwork.ParentNetworkId = &Reference{ Id: GetExternalNetwork(networkName).Id }
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
							Gateway: gatewayCidrArr[0],
							PrefixLength: prefixLen,
							DnsSuffix: "",
							DnsServer1: "",
							DnsServer2: "",
						},
					},
				}
			}

			if networkType == "NAT_ROUTED" {
				newVdcNetwork.Connection = &ConnectionInfo{
					RouterRef: Reference{ Id: GetEdge(gatewayName, orgvdcName).Id },
					ConnectionTypeValue: "INTERNAL",
				}
			}

			data, err := json.Marshal(newVdcNetwork)
			if err != nil {
				Fatal(err)
			}

			header := map[string]string{ "Content-Type": "application/json" }
			res := client.Request("POST", "/cloudapi/1.0.0/orgVdcNetworks", header, data)
			fmt.Println(string(res.Body))
		},
	}
	cmd.PersistentFlags().StringVarP(&orgvdcName, "orgvdc", "v", "", "org vdc name")
	cmd.PersistentFlags().StringVarP(&networkType, "type", "t", "", "network type (NAT_ROUTED | ISOLATED | DIRECT)")
	cmd.PersistentFlags().StringVarP(&gatewayName, "gateway", "g", "", "gateway name")
	cmd.PersistentFlags().StringVarP(&gatewayCidr, "cidr", "", "", "gateway cidr")
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
