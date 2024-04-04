package module

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
)

func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get resources or exec get api",
		Args:  cobra.MaximumNArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initClient()
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			api := args[0]
			if validateApi(api) {
				res := client.Request("GET", api, nil, nil)
				fmt.Println(string(res.Body))
			} else {
				Fatal("\"" + api + "\" is not a valid command or api")
			}
		},
	}
	cmd.AddCommand(
		NewCmdGetOrg(),
		NewCmdGetOrgVdc(),
		NewCmdGetOrgVdcNetwork(),
		NewCmdGetVApp(),
		NewCmdGetVAppNetwork(),
		NewCmdGetVAppVm(),
		NewCmdGetVAppVmNetwork(),
		NewCmdGetTask(),
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

func NewCmdGetOrgVdcNetwork() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vdc-network ${VDC_NAME}",
		Short:   "Get VdcNetwork [vn]",
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

func NewCmdGetVApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vapp",
		Aliases: []string{"a"},
		Short:   "Get VApp [a]",
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

			var data [][]string
			for _, vm := range GetVAppVm(vapp.Id) {
				vnics := vm.NetworkConnectionSection.NetworkConnection
				sort.Slice(vnics, func(i, j int) bool {
					return vnics[i].NetworkConnectionIndex < vnics[j].NetworkConnectionIndex
				})
				for _, nw := range vnics {
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

func NewCmdGetTask() *cobra.Command {
	var taskId string
	var latest bool

	cmd := &cobra.Command{
		Use:     "task",
		Aliases: []string{"t"},
		Short:   "Get Tasks of Org [t]",
		Args:    cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			initClient()
			orgs := GetOrgs()
			orgNames := []string{}
			for _, o := range orgs {
				orgNames = append(orgNames, o.Name)
			}

			return orgNames, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			orgId := ""
			orgName := client.site.OrgName
			if len(args) > 0 {
				orgName = args[0]
			}
			if orgName != "" {
				org := GetOrg(orgName)
				orgId = org.Id
			}

			if latest {
				task := GetTasks(orgId)[0]
				taskId = LastOne(task.Href, "/")
			}

			if taskId == "" {
				header := []string{"Org", "Operation", "Id", "Status", "Object", "User", "Start", "End"}
				maxcount := 5
				if orgId == "" {
					maxcount = 10
				}
				var data [][]string
				for _, task := range GetTasks(orgId) {
					data = append(data, []string{
						task.Org.Name,
						task.Operation,
						LastOne(task.Href, "/"),
						task.Status,
						task.Owner.Name,
						task.User.Name,
						task.StartTime,
						task.EndTime})
					maxcount--
					if maxcount == 0 {
						break
					}
				}
				PrityPrint(header, data)
			} else {
				task := GetTask(taskId)
				fmt.Println("Org: " + task.Org.Name)
				fmt.Println("Operation: " + task.Operation)
				fmt.Println("Status: " + task.Status)
				fmt.Println("Time: " + task.StartTime + " - " + task.EndTime)
				fmt.Println("Object: " + task.Owner.Name + " (" + task.Owner.Href + ")")
				fmt.Println("User: " + task.User.Name)
				if task.Error.TenantError.Message != "" {
					fmt.Println("Error: " + task.Error.TenantError.Message)
					fmt.Println("ErrorCode: " + task.Error.TenantError.MajorErrorCode + " - " + task.Error.TenantError.MinorErrorCode)
					fmt.Println("Trace: " + task.Error.StackTrace)
				}

				if len(task.VcTaskList.VcTask) > 0 {
					fmt.Println("VcTasks: ")
					header := []string{"Operation", "Status", "ObjectType", "ObjectName", "ObjectId", "Start", "End"}
					var data [][]string
					for _, vt := range task.VcTaskList.VcTask {
						data = append(data, []string{
							vt.Description,
							vt.Status,
							vt.ObjectType,
							vt.ObjectName,
							vt.ObjectMoref,
							vt.StartTime,
							vt.EndTime})
					}
					PrityPrint(header, data)
				}
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&taskId, "id", "i", "", "task id")
	cmd.PersistentFlags().BoolVarP(&latest, "latest", "l", false, "latest task")
	return cmd
}

func GetOrgs() []Org {
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

func GetOrgVdcs() []OrgVdc {
	var vdcs []OrgVdc
	var orgVdcList OrgVdcList

	res := client.Request("GET", "/api/query?type=adminOrgVdc", nil, nil)
	err := xml.Unmarshal(res.Body, &orgVdcList)
	if err != nil {
		Fatal(err)
	}
	vdcs = append(vdcs, orgVdcList.OrgVdc...)

	for orgVdcList.Total > orgVdcList.PageSize*orgVdcList.Page {
		api := fmt.Sprintf("/api/query?type=adminOrgVdc&pageSize=%d&page=%d", orgVdcList.PageSize, orgVdcList.Page+1)
		res := client.Request("GET", api, nil, nil)
		orgVdcList = OrgVdcList{}
		err := xml.Unmarshal(res.Body, &orgVdcList)
		if err != nil {
			Fatal(err)
		}
		vdcs = append(vdcs, orgVdcList.OrgVdc...)
	}

	for i := 0; i < len(vdcs); i++ {
		vdcs[i].Id = LastOne(vdcs[i].Href, "/")
	}

	sort.Slice(vdcs, func(i, j int) bool {
		return vdcs[i].Name < vdcs[j].Name
	})
	return vdcs
}

func GetOrgVdcNetwork(vdcId string) []OrgVdcNetwork {
	res := client.Request("GET", "/api/admin/vdc/"+vdcId, nil, nil)

	var adminVdc AdminVdc
	err := xml.Unmarshal(res.Body, &adminVdc)
	if err != nil {
		Fatal(err)
	}

	var orgVdcNetworkList []OrgVdcNetwork
	for i := 0; i < len(adminVdc.AvailableNetworks.Network); i++ {
		networkId := LastOne(adminVdc.AvailableNetworks.Network[i].Href, "/")
		res2 := client.Request("GET", "/api/admin/network/"+networkId, nil, nil)

		var orgVdcNetwork OrgVdcNetwork
		err := xml.Unmarshal(res2.Body, &orgVdcNetwork)
		if err != nil {
			Fatal(err)
		}
		orgVdcNetwork.Id = LastOne(orgVdcNetwork.Href, "/")
		orgVdcNetworkList = append(orgVdcNetworkList, orgVdcNetwork)
	}

	sort.Slice(orgVdcNetworkList, func(i, j int) bool {
		return orgVdcNetworkList[i].Name < orgVdcNetworkList[j].Name
	})

	return orgVdcNetworkList
}

func GetVApps() []VApp {
	vapps := []VApp{}
	page := 1
	for {
		res := client.Request("GET", fmt.Sprintf("/api/admin/extension/vapps/query?page=%d", page), nil, nil)
		var vappList VAppList
		if err := xml.Unmarshal(res.Body, &vappList); err != nil {
			Fatal(err)
		}
		vapps = append(vapps, vappList.VApp...)
		if vappList.Total <= vappList.Page * vappList.PageSize {
			break
		}
		page++
	}

	var orgList []Org = GetOrgs()

	for i := 0; i < len(vapps); i++ {
		vapps[i].Id = LastOne(vapps[i].Href, "/")
		for _, org := range orgList {
			if vapps[i].OrgHref == org.Href {
				vapps[i].OrgName = org.Name
				break
			}
		}
	}

	sort.Slice(vapps, func(i, j int) bool {
		return vapps[i].Name < vapps[j].Name
	})

	return vapps
}

func GetVdcNetworkType(networkId string) string {
	res := client.Request("GET", "/api/network/"+networkId, nil, nil)

	var network Network
	err := xml.Unmarshal(res.Body, &network)
	if err != nil {
		Fatal(err)
	}

	return network.Configuration.FenceMode
}

func GetTasks(orgId string) []Task {
	tasks := []Task{}

	orgIdList := []string{}
	if orgId != "" {
		orgIdList = append(orgIdList, orgId)
		system := GetOrg("System")
		orgIdList = append(orgIdList, system.Id)
	} else {
		for _, org := range GetOrgs() {
			orgIdList = append(orgIdList, org.Id)
		}
	}

	for _, id := range orgIdList {
		res := client.Request("GET", "/api/tasksList/"+id, nil, nil)

		var taskList TaskList
		err := xml.Unmarshal(res.Body, &taskList)
		if err != nil {
			Fatal(err)
		}

		tasks = append(tasks, taskList.Task...)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].StartTime > tasks[j].StartTime
	})

	return tasks
}

func GetTask(taskId string) Task {
	res := client.Request("GET", "/api/task/"+taskId, nil, nil)

	var task Task
	err := xml.Unmarshal(res.Body, &task)
	if err != nil {
		Fatal(err)
	}

	return task
}

func GetOrg(orgName string) Org {
	res := client.Request("GET", fmt.Sprintf("/api/admin/orgs/query?filter=(name==%s)", orgName), nil, nil)

	var orgResults struct {
		OrgRecord Org `xml:"OrgRecord"`
	}
	err := xml.Unmarshal(res.Body, &orgResults)
	if err != nil {
		Fatal(err)
	}
	orgResults.OrgRecord.Id = LastOne(orgResults.OrgRecord.Href, "/")

	return orgResults.OrgRecord
}

func GetVdc(vdcName string) (OrgVdc, error) {
	for _, vdc := range GetOrgVdcs() {
		if vdc.Name == vdcName {
			return vdc, nil
		}
	}
	return OrgVdc{}, fmt.Errorf("Org VDC \"" + vdcName + "\" not found")
}

func GetVAppByName(vappName string) (VApp, error) {
	for _, vapp := range GetVApps() {
		if vapp.Name == vappName {
			return vapp, nil
		}
	}
	return VApp{}, fmt.Errorf("vApp \"%s\" not found", vappName)
}

func GetVAppNetwork(vappId string) []Network {
	res := client.Request("GET", "/api/vApp/"+vappId+"/networkConfigSection", nil, nil)

	var networkConfigSection NetworkConfigSection
	err := xml.Unmarshal(res.Body, &networkConfigSection)
	if err != nil {
		Fatal(err)
	}

	nws := networkConfigSection.NetworkConfig
	sort.Slice(nws, func(i, j int) bool {
		return nws[i].Name < nws[j].Name
	})

	return nws
}

func GetVAppVm(vappId string) []VM {
	res := client.Request("GET", "/api/vApp/"+vappId, nil, nil)

	var vappDetails VAppDetails
	err := xml.Unmarshal(res.Body, &vappDetails)
	if err != nil {
		Fatal(err)
	}

	vms := vappDetails.VMs.VM
	sort.Slice(vms, func(i, j int) bool {
		return vms[i].Name < vms[j].Name
	})

	return vms
}

func GetProviderVdc(name string) VdcReference {
	res := client.Request("GET", fmt.Sprintf("/api/admin/extension/providerVdcReferences/query?filter=(name==%s)", name), nil, nil)

	result := struct {
		VMWProviderVdcRecord struct {
			Name string `xml:"name,attr"`
			Href string `xml:"href,attr"`
		} `xml:"VMWProviderVdcRecord"`
	}{}
	err := xml.Unmarshal(res.Body, &result)
	if err != nil {
		Fatal(err)
	}

	return VdcReference {
		Name: result.VMWProviderVdcRecord.Name,
		Href: result.VMWProviderVdcRecord.Href,
		Id: fmt.Sprintf("urn:vcloud:providervdc:%s", LastOne(result.VMWProviderVdcRecord.Href, "/")),
	}
}

func GetNetworkPool(name string) VdcReference {
	res := client.Request("GET", fmt.Sprintf("/api/admin/extension/networkPoolReferences/query?filter=(name==%s)", name), nil, nil)

	result := struct {
		NetworkPoolRecord struct {
			Name string `xml:"name,attr"`
			Href string `xml:"href,attr"`
		} `xml:"NetworkPoolRecord"`
	}{}
	err := xml.Unmarshal(res.Body, &result)
	if err != nil {
		Fatal(err)
	}

	return VdcReference {
		Name: result.NetworkPoolRecord.Name,
		Href: result.NetworkPoolRecord.Href,
		Id: fmt.Sprintf("urn:vcloud:networkpool:%s", LastOne(result.NetworkPoolRecord.Href, "/")),
	}
}

func GetStorageProfile(name string, providerVdcName string) VdcReference {
	api := fmt.Sprintf("/cloudapi/1.0.0/pvdcStoragePolicies?filter=(name==%s;providerVdcRef.name==%s)", name, providerVdcName)
	res := client.Request("GET", api, nil, nil)

	result := struct {
		Values []struct {
			Urn  string `json:"id"`
			Name string `json:"name"`
		} `json:"values"`
	}{}
	err := json.Unmarshal(res.Body, &result)
	if err != nil {
		Fatal(err)
	}

	return VdcReference {
		Name: result.Values[0].Name,
		Id: fmt.Sprintf("urn:vcloud:providervdcstorageprofile:%s", LastOne(result.Values[0].Urn, ":")),
	}
}
