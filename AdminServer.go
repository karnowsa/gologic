package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/go-resty/resty"
)

type AdminServer struct {
	ipAdress       string
	port           int
	username       string
	password       string
	weblogicHome   string
	middlewareHome string
	Cli            *resty.Client
	ManagedList    map[string]*ManagedServer
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
	admin.ManagedList = make(map[string]*ManagedServer, 30)
	admin.Cli = resty.New()
	admin.Cli.SetBasicAuth(admin.username, admin.password)
	admin.Cli.SetDisableWarn(true)
	admin.Cli.SetHostURL("http://" + admin.ipAdress + ":" + strconv.Itoa(admin.port))
	resp, err := admin.Cli.R().
		EnableTrace().
		SetHeader("Accept", "application/json").
		Get("/management/weblogic/latest/domainRuntime/serverLifeCycleRuntimes?links=none&fields=name,state,weblogicHome,middlewareHome")

	if err != nil {
		panic(err)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(fmt.Sprintf("%v", resp)), &result)

	items := result["items"].([]interface{})
	for _, value := range items {
		admin.ManagedList[value.(map[string]interface{})["name"].(string)] = &ManagedServer{
			Name:           value.(map[string]interface{})["name"].(string),
			Status:         value.(map[string]interface{})["state"].(string),
			WeblogicHome:   value.(map[string]interface{})["weblogicHome"].(string),
			MiddlewareHome: value.(map[string]interface{})["middlewareHome"].(string),
			Cli:            admin.Cli}
	}
}

func (admin *AdminServer) status(nameList []string) {
	fmt.Println()
	if nameList == nil {
		fmt.Printf("%-40s %-15s \n", "AdminServer", admin.ManagedList["AdminServer"].statusMS())
		fmt.Println()
		fmt.Printf("---------------------------------------------------------\n")
		fmt.Println()
		for _, name := range admin.sortedManagedList() {
			if name != "AdminServer" {
				fmt.Printf("%-40s %-15s \n", name, admin.ManagedList[name].statusMS())
			}
		}
	} else {
		for _, name := range nameList {
			managedserver, ok := admin.ManagedList[name]
			if ok {
				fmt.Printf("%-40s %-15s \n", name, managedserver.statusMS())
			}
		}

	}
	fmt.Println()

}

func (admin *AdminServer) start(nameList []string) {
	fmt.Println()
	if nameList == nil {
		for name := range admin.ManagedList {
			if name != "AdminServer" {
				admin.ManagedList[name].startMS()
				fmt.Printf("%-40s %-15s \n", name, admin.ManagedList[name].statusMS())
			}
		}
	} else {
		for _, name := range nameList {
			managedserver, ok := admin.ManagedList[name]
			if ok {
				managedserver.startMS()
				fmt.Printf("%-40s %-15s \n", name, managedserver.statusMS())
			}
		}
	}
	fmt.Println()
}

func (admin *AdminServer) stop(nameList []string) {
	if nameList == nil {
		for name := range admin.ManagedList {
			if name != "AdminServer" {
				admin.ManagedList[name].stopMS()
			}
		}
	} else {
		for _, name := range nameList {
			managedserver, ok := admin.ManagedList[name]
			if ok {
				managedserver.stopMS()
			}
		}
	}
}

func (admin *AdminServer) deploy() {
}

func (admin *AdminServer) createManagedServer() {
}
