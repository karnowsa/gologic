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
	cli            *resty.Client
	managedList    map[string]ManagedServer
}

func (admin *AdminServer) sortedManagedList() []string {
	keys := make([]string, 0, len(admin.managedList))
	for k := range admin.managedList {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func (admin *AdminServer) init() {
	admin.managedList = make(map[string]ManagedServer, 30)
	admin.cli = resty.New()
	admin.cli.SetBasicAuth(admin.username, admin.password)
	admin.cli.SetDisableWarn(true)
	admin.cli.SetHostURL("http://" + admin.ipAdress + ":" + strconv.Itoa(admin.port))
	resp, err := admin.cli.R().
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
		admin.managedList[value.(map[string]interface{})["name"].(string)] = ManagedServer{
			status:         value.(map[string]interface{})["state"].(string),
			weblogicHome:   value.(map[string]interface{})["weblogicHome"].(string),
			middlewareHome: value.(map[string]interface{})["middlewareHome"].(string),
			cli:            admin.cli}
	}
}

func (admin *AdminServer) status(nameList []string) {
	fmt.Printf("%-40s %-15s \n", "AdminServer", admin.managedList["AdminServer"].status)
	fmt.Printf("---------------------------------------------------------\n")

	if nameList == nil {
		for _, name := range admin.sortedManagedList() {
			if name != "AdminServer" {
				fmt.Printf("%-40s %-15s \n", name, admin.managedList[name].status)
			}
		}
	} else {
		for _, name := range nameList {
			value, ok := admin.managedList[name]
			if ok {
				fmt.Printf("%-40s %-15s \n", name, value.status)
			}
		}
	}
}

func (admin *AdminServer) start() {

}

func (admin *AdminServer) startAll() {
	for name := range admin.managedList {
		if name != "AdminServer" {
			resp, err := admin.cli.R().
				SetPathParams(map[string]string{
					"managedServerName": name,
				}).
				Get("/management/weblogic/latest/domainRuntime/serverLifeCycleRuntimes/{managedServerName}/start")

			if err != nil {
				panic(err)
			}

			fmt.Println(resp)
		}
	}

}

func (admin *AdminServer) stop() {
}

func (admin *AdminServer) stopAll() {
}

func (admin *AdminServer) deploy() {
}

func (admin *AdminServer) createManagedServer() {
}
