package module

type OrgList struct {
	Org []Org `xml:"Org"`
}

type Org struct {
	Id   string
	Href string `xml:"href,attr"`
	Name string `xml:"name,attr"`
}

type OrgVdcList struct {
	OrgVdc []OrgVdc `xml:"AdminVdcRecord"`
}

type OrgVdc struct {
	Name                  string `xml:"name,attr"`
	Href                  string `xml:"href,attr"`
	Id                    string
	IsEnabled             string `xml:"isEnabled,attr"`
	OrgName               string `xml:"orgName,attr"`
	ProviderVdcName       string `xml:"providerVdcName,attr"`
	VcName                string `xml:"vcName,attr"`
	NetworkProviderType   string `xml:"networkProviderType,attr"`
	NumberOfVApps         int    `xml:"numberOfVApps,attr"`
	NumberOfVMs           int    `xml:"numberOfVMs,attr"`
	NumberOfVAppTemplates int    `xml:"numberOfVAppTemplates,attr"`
}

type VAppList struct {
	VApp []VApp `xml:"AdminVAppRecord"`
}

type VApp struct {
	Name           string `xml:"name,attr"`
	Href           string `xml:"href,attr"`
	Id             string
	IsEnabled      string `xml:"isEnabled,attr"`
	Status         string `xml:"status,attr"`
	OrgHref        string `xml:"org,attr"`
	OrgName        string
	VdcName        string `xml:"vdcName,attr"`
	NumberOfVMs    int    `xml:"numberOfVMs,attr"`
	TaskStatusName string `xml:"taskStatusName,attr"`
	TaskStatus     string `xml:"taskStatus,attr"`
}

type OrgVdcNetworkList struct {
	OrgVdcNetworks []OrgVdcNetwork `xml:"OrgVdcNetworkRecord"`
}

type Element struct {
	Href string `xml:"href,attr"`
	Id   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type OrgVdcNetwork struct {
	Name               string `xml:"name,attr"`
	Href               string `xml:"href,attr"`
	Id                 string
	DefaultGateway     string `xml:"defaultGateway,attr"`
	Netmask            string `xml:"netmask,attr"`
	SubnetPrefixLength string `xml:"subnetPrefixLength,attr"`
	Dns1               string `xml:"dns1,attr"`
	Dns2               string `xml:"dns2,attr"`
	DnsSuffix          string `xml:"dnsSuffix,attr"`
	OrgName            string
	VdcHref            string `xml:"vdc,attr"`
	VdcName            string `xml:"vdcName,attr"`
	IsShared           string `xml:"isShared,attr"`
	IsIpScopeInherited string `xml:"isIpScopeInherited,attr"`
	FenceMode          string
}

type NetworkConfigSection struct {
	NetworkConfig []Network `xml:"NetworkConfig"`
}

type Network struct {
	Name          string               `xml:"networkName,attr"`
	Configuration NetworkConfiguration `xml:"Configuration"`
}

type NetworkConfiguration struct {
	IpScopes      IpScopeList `xml:"IpScopes"`
	ParentNetwork Element     `xml:"ParentNetwork,omitempty"`
	FenceMode     string      `xml:"FenceMode"`
}

type IpScopeList struct {
	IpScope []IpScope `xml:"IpScope"`
}

type IpScope struct {
	IsInherited        string `xml:"IsInherited"`
	Gateway            string `xml:"Gateway"`
	Netmask            string `xml:"Netmask"`
	SubnetPrefixLength string `xml:"SubnetPrefixLength"`
	IsEnabled          string `xml:"IsEnabled"`
}

type VAppDetails struct {
	VMs VmList `xml:"Children"`
}

type VmList struct {
	VM []VM `xml:"Vm"`
}

type VM struct {
	Name                     string                   `xml:"name,attr"`
	Urn                      string                   `xml:"id,attr"`
	Href                     string                   `xml:"href,attr"`
	NetworkConnectionSection NetworkConnectionSection `xml:"NetworkConnectionSection"`
}

type NetworkConnectionSection struct {
	PrimaryNetworkConnectionIndex int                 `xml:"PrimaryNetworkConnectionIndex"`
	NetworkConnection             []NetworkConnection `xml:"NetworkConnection"`
}

type NetworkConnection struct {
	Name                             string `xml:"network,attr"`
	NetworkConnectionIndex           int    `xml:"NetworkConnectionIndex"`
	IpAddress                        string `xml:"IpAddress"`
	IpType                           string `xml:"IpType"`
	IsConnected                      string `xml:"IsConnected"`
	MACAddress                       string `xml:"MACAddress"`
	IpAddressAllocationMode          string `xml:"IpAddressAllocationMode"`
	SecondaryIpAddressAllocationMode string `xml:"SecondaryIpAddressAllocationMode"`
	NetworkAdapterType               string `xml:"NetworkAdapterType"`
}
