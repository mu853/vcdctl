package module

type OrgList struct {
	Org []Org `xml:"Org"`
}

type Org struct {
	Id   string
	Href string `xml:"href,attr"`
	Name string `xml:"name,attr"`
}

type OrgJson struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Id          string `json:"id"`
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
	Enabled                   bool      `xml:"Enabled"`
	Units                     string    `xml:"Units"`
	Limit                     int       `xml:"Limit"`
	Default                   bool      `xml:"Default"`
	ProviderVdcStorageProfile Reference `xml:"ProviderVdcStorageProfile"`
}

type ConnectionInfo struct {
	RouterRef           ReferenceJson `json:"routerRef"`
	ConnectionTypeValue string        `json:"connectionTypeValue"` //INTERNAL or NON_DISTRIBUTED
	Connected           bool          `json:"connected"`
}

type Subnet struct {
	Gateway      string   `json:"gateway"`
	PrefixLength int      `json:"prefixLength"`
	DnsSuffix    string   `json:"dnsSuffix"`
	DnsServer1   string   `json:"dnsServer1"`
	DnsServer2   string   `json:"dnsServer2"`
	IpRanges     *IpRange `json:"ipRanges,omitempty"`
	Enabled      bool     `json:"enabled,omitempty"`
}

type Subnets struct {
	Values []Subnet `json:"values"`
}

type VAppList struct {
	Page     int    `xml:"page,attr"`
	PageSize int    `xml:"pageSize,attr"`
	Total    int    `xml:"total,attr"`
	VApp     []VApp `xml:"VAppRecord"`
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
	Network []Reference `xml:"Network"`
}

type Reference struct {
	Name string `xml:"name,attr"`
	Href string `xml:"href,attr"`
	Id   string `xml:"id,attr"`
}

type ReferenceJson struct {
	Name string `json:"name,omitempty"`
	Urn  string `json:"id"`
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
	ParentNetwork        *Reference  `xml:"ParentNetwork,omitempty"`
	FenceMode            string      `xml:"FenceMode"`
	DistributedInterface string      `xml:"DistributedInterface,omitempty"`
	ServiceInterface     string      `xml:"ServiceInterface,omitempty"`
	GuestVlanAllowed     string      `xml:"GuestVlanAllowed,omitempty"`
	Connected            string      `xml:"Connected,omitempty"`
}

type OrgVdcNetworkJson struct {
	Urn                     string          `json:"id,omitempty"`
	Name                    string          `json:"name"`
	Subnets                 *Subnets        `json:"subnets,omitempty"`
	BackingNetworkId        string          `json:"backingNetworkId,omitempty"`
	BackingNetworkType      string          `json:"backingNetworkType,omitempty"`
	ParentNetworkId         *ReferenceJson  `json:"parentNetworkId,omitempty"`
	NetworkType             string          `json:"networkType"`
	OrgVdc                  *ReferenceJson  `json:"orgVdc,omitempty"`
	OwnerRef                ReferenceJson   `json:"ownerRef"`
	OrgRef                  *ReferenceJson  `json:"orgRef,omitempty"`
	OrgVdcIsNsxTBacked      bool            `json:"orgVdcIsNsxTBacked,omitempty"`
	Connection              *ConnectionInfo `json:"connection,omitempty"`
	IsDefaultNetwork        bool            `json:"isDefaultNetwork,omitempty"`
	Shared                  bool            `json:"shared,omitempty"`
	EnableDualSubnetNetwork bool            `json:"enableDualSubnetNetwork,omitempty"`
	GuestVlanTaggingAllowed bool            `json:"guestVlanTaggingAllowed,omitempty"`
	RetainNicResources      bool            `json:"retainNicResources,omitempty"`
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

type NetworkBacking struct {
	Name            string        `json:"name"`
	Id              string        `json:"backingId"`
	BackingType     string        `json:"backingType"`
	NetworkProvider ReferenceJson `json:"networkProvider"`
}

type NetworkBackings struct {
	Values []NetworkBacking `json:"values"`
}

type ProviderGateway struct {
	Urn             string                   `json:"id"`
	Name            string                   `json:"name"`
	Subnets         EdgeGatewayUplinkSubnets `json:"subnets"`
	NetworkBackings NetworkBackings          `json:"networkBackings"`
}

type EdgeGatewayUplinkSubnets struct {
	Values []EdgeGatewayUplinkSubnet `json:"values"`
}

type EdgeGatewayUplinkSubnet struct {
	GatewayAddress       string    `json:"gateway"`
	PrefixLength         int       `json:"prefixLength"`
	DnsSuffix            string    `json:"dnsSuffix"`
	DnsServer1           string    `json:"dnsServer1"`
	DnsServer2           string    `json:"dnsServer2"`
	IpRanges             *IpRanges `json:"ipRanges,omitempty"`
	PrimaryIp            string    `json:"primaryIp,omitempty"`
	AutoAllocateIpRanges bool      `json:"autoAllocateIpRanges"`
	TotalIpCount         int       `json:"totalIpCount,omitempty"`
}

type EdgeGatewayUplink struct {
	UplinkId      string                   `json:"uplinkId"`
	UplinkName    string                   `json:"uplinkName"`
	Subnets       EdgeGatewayUplinkSubnets `json:"subnets"`
	Dedicated     bool                     `json:"dedicated"`
	Connected     bool                     `json:"connected"`
	UsingIpSpace  bool                     `json:"usingIpSpace"`
	VrfLiteBacked bool                     `json:"vrfLiteBacked"`
	BackingType   string                   `json:"backingType"`
}

type EdgeGateway struct {
	Urn                                      string              `json:"id,omitempty"`
	Name                                     string              `json:"name"`
	EdgeGatewayUplinks                       []EdgeGatewayUplink `json:"edgeGatewayUplinks"`
	DistributedRoutingEnabled                bool                `json:"distributedRoutingEnabled"`
	NonDistributedRoutingEnabled             bool                `json:"nonDistributedRoutingEnabled"`
	ServiceNetworkDefinition                 string              `json:"serviceNetworkDefinition"`
	DistributedRouterUplinkNetworkDefinition string              `json:"distributedRouterUplinkNetworkDefinition"`
	OrgVdcNetworkCount                       int                 `json:"orgVdcNetworkCount,omitempty"`
	GatewayBacking                           *NetworkBacking     `json:"gatewayBacking,omitempty"`
	EdgeClusterConfig                        *EdgeClusterConfig  `json:"edgeClusterConfig,omitempty"`
	OwnerRef                                 ReferenceJson       `json:"ownerRef"`
	OrgVdc                                   *ReferenceJson      `json:"orgVdc,omitempty"`
	OrgRef                                   *ReferenceJson      `json:"orgRef,omitempty"`
}

type EdgeClusterConfig struct {
	PrimaryEdgeCluster   EdgeClusterConfigSub  `json:"primaryEdgeCluster"`
	SecondaryEdgeCluster *EdgeClusterConfigSub `json:"secondaryEdgeCluster,omitempty"`
}

type EdgeClusterConfigSub struct {
	EdgeClusterRef ReferenceJson `json:"edgeClusterRef"`
	BackingId      string        `json:"backingId"`
}

type IpRanges struct {
	Values []*IpRange `json:"values,omitempty"`
}

type IpRange struct {
	StartAddress string `json:"startAddress,omitempty"`
	EndAddress   string `json:"endAddress,omitempty"`
}

type LeaseSettingsSection struct {
	Xmlns                     string `xml:"xmlns,attr"`
	XmlnsVmext                string `xml:"xmlns vmext,attr"`
	XmlnsOvf                  string `xml:"xmlns ovf,attr"`
	XmlnsVssd                 string `xml:"xmlns vssd,attr"`
	XmlnsCommon               string `xml:"xmlns common,attr"`
	XmlnsRasd                 string `xml:"xmlns rasd,attr"`
	XmlnsVmw                  string `xml:"xmlns vmw,attr"`
	XmlnsOvfenv               string `xml:"xmlns ovfenv,attr"`
	XmlnsNs9                  string `xml:"xmlns ns9,attr"`
	Href                      string `xml:"href,attr"`
	Type                      string `xml:"type,attr"`
	OvfRequired               bool   `xml:"ovf required,attr"`
	OvfInfo                   string `xml:"ovf Info"`

	DeploymentLeaseInSeconds  string `xml:"DeploymentLeaseInSeconds"`
	StorageLeaseInSeconds     string `xml:"StorageLeaseInSeconds"`
	DeploymentLeaseExpiration string `xml:"DeploymentLeaseExpiration,omitempty"`
	StorageLeaseExpiration    string `xml:"StorageLeaseExpiration,omitempty"`
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
	Operation     string      `xml:"operation,attr"`
	OperationName string      `xml:"operationName,attr"`
	Status        string      `xml:"status,attr"`
	StartTime     string      `xml:"startTime,attr"`
	EndTime       string      `xml:"endTime,attr"`
	Href          string      `xml:"href,attr"`
	Urn           string      `xml:"id,attr"`
	Org           Reference   `xml:"Organization"`
	User          Reference   `xml:"User"`
	Owner         Reference   `xml:"Owner"`
	Error         *TaskError  `xml:"Error,omitempty"`
	VcTaskList    *VcTaskList `xml:"VcTaskList,omitempty"`
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
