package gologic

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty"
)

/*ManagedServer is a a struct to save Weblogic ManagedServer
Name: It has a name
Status: And a status with color (Running, Shutdown...)
*/
type ManagedServer struct {
	Name   string
	Status string
	Cli    *resty.Client
}

//GetStatus returns the status of the ManagedServer with colors
func (ms *ManagedServer) GetStatus() string {
	if ms.Status == "RUNNING" {
		return "\033[32m[" + ms.Status + "]\033[0m"
	} else if ms.Status == "SHUTDOWN" {
		return "\033[31m[" + ms.Status + "]\033[0m"
	} else if ms.Status == "TASK IN PROGRESS" {
		return "\033[33m[" + ms.Status + "]\033[0m"
	} else if ms.Status == "StartING" {
		return "\033[36m[" + ms.Status + "]\033[0m"
	}
	return "\033[33m[" + ms.Status + "]\033[0m"
}

//StartMS starts a list of ManagedServer, when its empty then its starts all ManagedServer
func (ms *ManagedServer) StartMS() {
	var result map[string]interface{}

	resp, err := ms.Cli.R().
		SetPathParams(map[string]string{
			"managedServerName": ms.Name,
		}).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		SetHeader("Prefer", "respond-async").
		SetBody("{}").
		Post("/domainRuntime/serverLifeCycleRuntimes/{managedServerName}/start")

	if err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(fmt.Sprintf("%v", resp)), &result)

	taskStatus, ok := result["taskStatus"].(string)

	if ok {
		ms.Status = taskStatus
	} else {
		statusCode, ok := result["status"].(float64)
		if ok {
			if statusCode == 400 {
				ms.Status = "RUNNING"
			} else {
				panic(statusCode)
			}
		} else {
			panic(ok)
		}
	}

}

//StopMS stops a list of ManagedServer, when its empty then its Stops all ManagedServer
func (ms *ManagedServer) StopMS() {
	var result map[string]interface{}

	resp, err := ms.Cli.R().
		SetPathParams(map[string]string{
			"managedServerName": ms.Name,
		}).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		SetHeader("Prefer", "respond-async").
		SetBody("{}").
		Post("/domainRuntime/serverLifeCycleRuntimes/{managedServerName}/forceShutdown")

	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(fmt.Sprintf("%v", resp)), &result)

	taskStatus, ok := result["taskStatus"].(string)

	if ok {
		ms.Status = taskStatus
	} else {
		statusCode, ok := result["status"].(float64)
		if ok {
			if statusCode == 400 {
				ms.Status = "SHUTDOWN"
			} else {
				panic(statusCode)
			}
		} else {
			panic(ok)
		}
	}
}
