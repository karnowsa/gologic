package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/go-resty/resty"
)

type AdminServer struct {
	name        string
	ipAdress    string
	port        int
	username    string
	password    string
	statusAdmin string
	Cli         *resty.Client
	ManagedList map[string]*ManagedServer
}

func (admin *AdminServer) sortedManagedList() []string {
	keys := make([]string, 0, len(admin.ManagedList))
	for k := range admin.ManagedList {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func (admin *AdminServer) init() {
	var resp *resty.Response
	var err error
	var result map[string]interface{}

	admin.ManagedList = make(map[string]*ManagedServer, 30)
	admin.Cli = resty.New()
	admin.Cli.SetBasicAuth(admin.username, admin.password)
	admin.Cli.SetDisableWarn(true)
	admin.Cli.SetHostURL("http://" + admin.ipAdress + ":" + strconv.Itoa(admin.port) + "/management/weblogic/latest")

	admin.checkAdminStatus()

	resp, err = admin.Cli.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		Get("/edit")

	if err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(fmt.Sprintf("%v", resp)), &result)

	admin.name = result["adminServerName"].(string)

	resp, err = admin.Cli.R().
		EnableTrace().
		SetHeader("Accept", "application/json").
		Get("/domainRuntime/serverLifeCycleRuntimes?links=none&fields=name,state")

	if err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(fmt.Sprintf("%v", resp)), &result)

	items := result["items"].([]interface{})

	for _, value := range items {
		if value.(map[string]interface{})["name"].(string) != admin.name {
			admin.ManagedList[value.(map[string]interface{})["name"].(string)] = &ManagedServer{
				Name:   value.(map[string]interface{})["name"].(string),
				Status: value.(map[string]interface{})["state"].(string),
				Cli:    admin.Cli}
		}
	}
}

func (admin *AdminServer) checkAdminStatus() {
	_, err := admin.Cli.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		Get("/")

	if err != nil {
		panic(err)
	} else {
		admin.statusAdmin = "RUNNING"
	}
}

func (admin *AdminServer) getStatus() string {
	if admin.statusAdmin == "RUNNING" {
		return "\033[32m[" + admin.statusAdmin + "]\033[0m"
	} else if admin.statusAdmin == "SHUTDOWN" {
		return "\033[31m[" + admin.statusAdmin + "]\033[0m"
	}
	return "\033[33m[" + admin.statusAdmin + "]\033[0m"
}

func (admin *AdminServer) start(nameList []string) {
	fmt.Println()
	if nameList == nil {
		for name := range admin.ManagedList {
			if name != "AdminServer" {
				admin.ManagedList[name].startMS()
				fmt.Printf("%-40s %s \n", name, admin.ManagedList[name].getStatus())
			}
		}
	} else {
		for _, name := range nameList {
			managedserver, ok := admin.ManagedList[name]
			if ok {
				managedserver.startMS()
				fmt.Printf("%-40s %s \n", name, managedserver.getStatus())
			} else {
				fmt.Printf("%-40s %s \n", name, "Doesn't exists")
			}
		}
	}
	fmt.Println()
}

func (admin *AdminServer) stop(nameList []string) {
	fmt.Println()
	if nameList == nil {
		for name := range admin.ManagedList {
			if name != "AdminServer" {
				admin.ManagedList[name].stopMS()
				fmt.Printf("%-40s %s \n", name, admin.ManagedList[name].getStatus())
			}
		}
	} else {
		for _, name := range nameList {
			managedserver, ok := admin.ManagedList[name]
			if ok {
				managedserver.stopMS()
				fmt.Printf("%-40s %s \n", name, managedserver.getStatus())
			} else {
				fmt.Printf("%-40s %s \n", name, "Doesn't exists")
			}
		}
	}
	fmt.Println()
}

//A list of Managed Server can be passed to this receiver
func (admin *AdminServer) printStatus(nameList []string) {
	fmt.Println()
	if nameList == nil {
		fmt.Printf("%-40s %-15s \n", admin.name, admin.getStatus())
		fmt.Println()
		fmt.Printf("---------------------------------------------------------\n")
		fmt.Println()
		for _, name := range admin.sortedManagedList() {
			if name != "AdminServer" {
				fmt.Printf("%-40s %s \n", name, admin.ManagedList[name].getStatus())
			}
		}
	} else {
		for _, name := range nameList {
			managedserver, ok := admin.ManagedList[name]
			if ok {
				fmt.Printf("%-40s %s \n", name, managedserver.getStatus())
			}
		}

	}
	fmt.Println()

}

func (admin *AdminServer) printInfo() {
	var result map[string]interface{}
	var resp *resty.Response

	resp, err := admin.Cli.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		Get("/domainConfig?links=none&fields=name,rootDirectory")
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(fmt.Sprintf("%v", resp)), &result)

	fmt.Printf("%-40s %s \n", "DomainName", result["name"].(string))
	fmt.Printf("%-40s %s \n", "DomainHome", result["rootDirectory"].(string))

	resp, err = admin.Cli.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		Get("/serverRuntime/JVMRuntime?links=none&fields=javaVersion,OSName,OSVersion")
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(fmt.Sprintf("%v", resp)), &result)

	fmt.Printf("%-40s %s\n", "Java Version", result["javaVersion"].(string))
	fmt.Printf("%-40s %s %s\n", "OS Version", result["OSName"].(string), result["OSVersion"].(string))
}

func (admin *AdminServer) createManagedServer() {
}
