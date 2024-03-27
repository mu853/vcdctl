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
	PageSize int      `xml:"pageSize,attr"`
	Page     int      `xml:"page,attr"`
	Total    int      `xml:"total,attr"`
	OrgVdc   []OrgVdc `xml:"AdminVdcRecord"`
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
	Page     int    `xml:"page,attr"`
	PageSize int    `xml:"pageSize,attr"`
	Total    int    `xml:"total,attr"`
	VApp     []VApp `xml:"AdminVAppRecord"`
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
	//OrgVdcNetworks []OrgVdcNetwork `xml:"OrgVdcNetworkRecord"`
	OrgVdcNetworks []OrgVdcAvailableNetwork `xml:"AvailableNetworks"`
}

type AdminVdc struct {
	AvailableNetworks OrgVdcAvailableNetwork `xml:"AvailableNetworks"`
}

type OrgVdcAvailableNetwork struct {
	Network []Element `xml:"Network"`
}

type Element struct {
	Href string `xml:"href,attr"`
	Id   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type OrgVdcNetwork struct {
	Name          string               `xml:"name,attr"`
	Href          string               `xml:"href,attr"`
	Urn           string               `xml:"id,attr"`
	Configuration NetworkConfiguration `xml:"Configuration"`
	IsShared      string               `xml:"IsShared"`
	Id            string
}

type NetworkConfigSection struct {
	NetworkConfig []Network `xml:"NetworkConfig"`
}

type Network struct {
	Name          string               `xml:"networkName,attr"`
	Configuration NetworkConfiguration `xml:"Configuration"`
}

type NetworkConfiguration struct {
	IpScopes             IpScopeList `xml:"IpScopes"`
	ParentNetwork        Element     `xml:"ParentNetwork,omitempty"`
	FenceMode            string      `xml:"FenceMode"`
	DistributedInterface string      `xml:"DistributedInterface,omitempty"`
	ServiceInterface     string      `xml:"ServiceInterface,omitempty"`
	GuestVlanAllowed     string      `xml:"GuestVlanAllowed,omitempty"`
	Connected            string      `xml:"Connected,omitempty"`
}

type IpScopeList struct {
	IpScope []IpScope `xml:"IpScope"`
}

type IpScope struct {
	IsInherited        string `xml:"IsInherited"`
	Gateway            string `xml:"Gateway"`
	Netmask            string `xml:"Netmask"`
	SubnetPrefixLength string `xml:"SubnetPrefixLength"`
	Dns1               string `xml:"Dns1,omitempty"`
	Dns2               string `xml:"Dns2,omitempty"`
	DnsSuffix          string `xml:"DnsSuffix,omitempty"`
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

type TaskList struct {
	Task []Task `xml:"Task"`
}

type Task struct {
	Operation     string     `xml:"operation,attr"`
	OperationName string     `xml:"operationName,attr"`
	Status        string     `xml:"status,attr"`
	StartTime     string     `xml:"startTime,attr"`
	EndTime       string     `xml:"endTime,attr"`
	Href          string     `xml:"href,attr"`
	Urn           string     `xml:"id,attr"`
	Org           Element    `xml:"Organization"`
	User          Element    `xml:"User"`
	Owner         Element    `xml:"Owner"`
	Error         TaskError  `xml:"Error,omitempty"`
	VcTaskList    VcTaskList `xml:"VcTaskList,omitempty"`
}

type TaskError struct {
	StackTrace     string          `xml:"stackTrace,attr"`
	MajorErrorCode string          `xml:"majorErrorCode,attr"`
	MinorErrorCode string          `xml:"minorErrorCode,attr"`
	TenantError    TaskTenantError `xml:"TenantError"`
}

type TaskTenantError struct {
	Message        string `xml:"message,attr"`
	MajorErrorCode string `xml:"majorErrorCode,attr"`
	MinorErrorCode string `xml:"minorErrorCode,attr"`
}

type VcTaskList struct {
	VcTask []VcTask `xml:"VcTask"`
}

type VcTask struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description,attr"`
	Status      string `xml:"status,attr"`
	ObjectType  string `xml:"objectType,attr"`
	ObjectName  string `xml:"objectName,attr"`
	ObjectMoref string `xml:"objectMoref,attr"`
	StartTime   string `xml:"startTime,attr"`
	EndTime     string `xml:"endTime,attr"`
}
