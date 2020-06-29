package gologic

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/go-resty/resty"
)

/*
AdminServer is a struct, which represents the Weblogic AdminServer.
name: is the Admin Server Name
ipAdress: is the IPv4 Address of the Server
username: username of a Administration Account
password: its the password of the account

*/
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

//Init checks the connection to the AdminServer and it collect a list of ManagedServer from the AdminServer
func Init(ip string, port int, username string, password string) AdminServer {
	var resp *resty.Response
	var err error
	var result map[string]interface{}
	var admin AdminServer

	admin.ipAdress = ip
	admin.port = port
	admin.username = username
	admin.password = password

	//Create a list for the ManagedServer
	admin.ManagedList = make(map[string]*ManagedServer, 30)

	//Configure the REST Client
	admin.Cli = resty.New()
	admin.Cli.SetBasicAuth(admin.username, admin.password)
	admin.Cli.SetDisableWarn(true)
	admin.Cli.SetHostURL("http://" + admin.ipAdress + ":" + strconv.Itoa(admin.port) + "/management/weblogic/latest")

	//Check the status of the AdminServer
	admin.checkAdminStatus()

	if resp, err = admin.Cli.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		Get("/edit"); err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(fmt.Sprintf("%v", resp)), &result)

	admin.name = result["adminServerName"].(string)

	if resp, err = admin.Cli.R().
		SetHeader("Accept", "application/json").
		Get("/domainRuntime/serverLifeCycleRuntimes?links=none&fields=name,state"); err != nil {
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

	return admin
}

//checkAdminStatus checks the status of the AdminServer
func (admin *AdminServer) checkAdminStatus() {
	if _, err := admin.Cli.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		Get("/"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		admin.statusAdmin = "RUNNING"
	}
}

//GetStatus returns the AdminServer Status with color
func (admin *AdminServer) GetStatus() string {
	if admin.statusAdmin == "RUNNING" {
		return "\033[32m[" + admin.statusAdmin + "]\033[0m"
	} else if admin.statusAdmin == "SHUTDOWN" {
		return "\033[31m[" + admin.statusAdmin + "]\033[0m"
	}
	return "\033[33m[" + admin.statusAdmin + "]\033[0m"
}

//Start starts a list of ManagedServer or when the list is empty, its stops every ManagedServer
func (admin *AdminServer) Start(nameList []string) {
	var err error
	fmt.Println()
	if len(nameList) <= 0 {
		for name := range admin.ManagedList {
			if name != "AdminServer" {
				if err = admin.ManagedList[name].StartMS(); err != nil {
					panic(err)
				}
				fmt.Printf("%-40v %v \n", name, admin.ManagedList[name].GetStatus())
			}
		}
	} else {
		for _, name := range nameList {
			managedserver, ok := admin.ManagedList[name]
			if ok {
				managedserver.StartMS()
				fmt.Printf("%-40v %v \n", name, managedserver.GetStatus())
			} else {
				fmt.Printf("%-40v %v \n", name, "Doesn't exists")
			}
		}
	}
	fmt.Println()
}

//Stop stops a list of servers or when the list is empty, its stops every ManagedServer
func (admin *AdminServer) Stop(nameList []string) {
	var err error

	fmt.Println()
	if len(nameList) <= 0 {
		for name := range admin.ManagedList {
			if name != "AdminServer" {
				if err = admin.ManagedList[name].StopMS(); err != nil {
					panic(err)
				}
				fmt.Printf("%-40v %v \n", name, admin.ManagedList[name].GetStatus())
			}
		}
	} else {
		for _, name := range nameList {
			managedserver, ok := admin.ManagedList[name]
			if ok {
				if err = managedserver.StopMS(); err != nil {
					panic(err)
				}
				fmt.Printf("%-40v %v \n", name, managedserver.GetStatus())
			} else {
				fmt.Printf("%-40v %v \n", name, "Doesn't exists")
			}
		}
	}
	fmt.Println()
}

//PrintStatus prints the status of all Servers or a list of specific servers
func (admin *AdminServer) PrintStatus(nameList []string) {
	fmt.Println()
	if len(nameList) <= 0 {
		fmt.Printf("%-40s %-15s \n", admin.name, admin.GetStatus())
		fmt.Println()
		fmt.Printf("---------------------------------------------------------\n")
		fmt.Println()
		for _, name := range admin.sortedManagedList() {
			if name != "AdminServer" {
				fmt.Printf("%-40v %v \n", name, admin.ManagedList[name].GetStatus())
			}
		}
	} else {
		for _, name := range nameList {
			managedserver, ok := admin.ManagedList[name]
			if ok {
				fmt.Printf("%-40v %v \n", name, managedserver.GetStatus())
			} else if name == admin.name {
				fmt.Printf("%-40v %v \n", name, admin.GetStatus())
			}
		}

	}
	fmt.Println()

}

//PrintInfo prints informations about the AdminServer
func (admin *AdminServer) PrintInfo() {
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

	fmt.Printf("%-40v %v \n", "AdminServer", admin.ipAdress+":"+strconv.Itoa(admin.port))
	fmt.Printf("%-40v %v \n", "DomainName", result["name"])
	fmt.Printf("%-40v %v \n", "DomainHome", result["rootDirectory"])

	resp, err = admin.Cli.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		Get("/serverRuntime/JVMRuntime?links=none&fields=javaVersion,OSName,OSVersion")
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(fmt.Sprintf("%v", resp)), &result)

	fmt.Printf("%-40v %v\n", "JavaVersion", result["javaVersion"])
	fmt.Printf("%-40v %v %v\n", "OSVersion", result["OSName"], result["OSVersion"])
}

//PrintDeployments prints a list of Deployments
func (admin *AdminServer) PrintDeployments() {
	var result map[string]interface{}
	var resp *resty.Response
	var err error

	resp, err = admin.Cli.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		Get("/domainRuntime/deploymentManager/appDeploymentRuntimes?links=none")

	if err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(fmt.Sprintf("%v", resp)), &result)

	items := result["items"].([]interface{})
	if len(items) > 0 {
		for _, value := range items {
			applicationEntry := value.(map[string]interface{})

			fmt.Printf("%-20s %-10v \n", applicationEntry["name"], applicationEntry["applicationVersion"])
		}

	} else {
		fmt.Println("No Deployments found")
	}
}

//CreateManagedServer creates a ManagedServer with the parameter name (name of the server), listenAddress, listenPort
func (admin *AdminServer) CreateManagedServer(name string, listenAddress string, listenPort string) {
	var result map[string]interface{}
	var resp *resty.Response
	var err error

	//Get a form to create a managed server
	if resp, err = admin.Cli.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		SetBody("{}").
		Post("/edit/changeManager/StartEdit"); err != nil {
		panic(err)
	}

	//Get a form to create a managed server
	if resp, err = admin.Cli.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		Get("/edit/serverCreateForm?links=none"); err != nil {
		panic(err)
	}

	fmt.Printf("Creating server (%v, %v, %v)\n", name, listenAddress, listenPort)

	if resp, err = admin.Cli.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		SetBody(`{ "name":"` + name + `", "listenAddress":"` + listenAddress + `", "listenPort":` + listenPort + `}`).
		Post("/edit/servers"); err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(fmt.Sprintf("%v", resp)), &result)

	if fmt.Sprintf("%v", resp) != "{}" {
		if strings.Contains(result["detail"].(string), "already exists") {
			fmt.Printf("The Managed Server %v already exists!\n", name)
		} else {
			fmt.Printf("Couldn't create Server %v!\n", name)
		}
	}

	if resp, err = admin.Cli.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		SetBody("{}").
		Post("/edit/changeManager/activate"); err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(fmt.Sprintf("%v", resp)), &result)

	if fmt.Sprintf("%v", resp) != "{}" {
		fmt.Println("Couldn't activate changes!")
	}

}
