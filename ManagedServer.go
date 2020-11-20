package github.com/karnowsa/gologic

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
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
func (ms *ManagedServer) StartMS() error {
	var result map[string]interface{}
	var resp *resty.Response
	var err error

	if resp, err = ms.Cli.R().
		SetPathParams(map[string]string{
			"managedServerName": ms.Name,
		}).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		SetHeader("Prefer", "respond-async").
		SetBody("{}").
		Post("/domainRuntime/serverLifeCycleRuntimes/{managedServerName}/start"); err != nil {
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
				return fmt.Errorf("StartMS() HTTP Statuscode is %v", statusCode)
			}
		} else {
			return fmt.Errorf("StartMS() ok is %v", ok)
		}
	}
	return nil
}

//StopMS stops a list of ManagedServer, when its empty then its Stops all ManagedServer
func (ms *ManagedServer) StopMS() error {
	var result map[string]interface{}
	var resp *resty.Response
	var err error

	if resp, err = ms.Cli.R().
		SetPathParams(map[string]string{
			"managedServerName": ms.Name,
		}).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		SetHeader("Prefer", "respond-async").
		SetBody("{}").
		Post("/domainRuntime/serverLifeCycleRuntimes/{managedServerName}/forceShutdown"); err != nil {
		return err
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
				return fmt.Errorf("StopMS() HTTP Statuscode is %v", statusCode)
			}
		} else {
			return fmt.Errorf("StopMS() ok is %v", ok)
		}
	}
	return nil
}
