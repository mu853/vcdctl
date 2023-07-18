package module

type OrgList struct {
	Orgs []Org `xml:"Org"`
}

type Org struct {
	Id   string
	Href string `xml:"href,attr"`
	Name string `xml:"name,attr"`
}

type OrgVdcList struct {
	OrgVdcs []OrgVdc `xml:"AdminVdcRecord"`
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
	VApps []VApp `xml:"AdminVAppRecord"`
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
}
