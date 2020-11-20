package github.com/karnowsa/gologic

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type ManagedServer struct {
	Name   string
	Status string
	Cli    *resty.Client
}

func (ms *ManagedServer) getStatus() string {
	if ms.Status == "RUNNING" {
		return "\033[32m[" + ms.Status + "]\033[0m"
	} else if ms.Status == "SHUTDOWN" {
		return "\033[31m[" + ms.Status + "]\033[0m"
	} else if ms.Status == "TASK IN PROGRESS" {
		return "\033[33m[" + ms.Status + "]\033[0m"
	} else if ms.Status == "STARTING" {
		return "\033[36m[" + ms.Status + "]\033[0m"
	}
	return "\033[33m[" + ms.Status + "]\033[0m"
}

func (ms *ManagedServer) startMS() {
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

func (ms *ManagedServer) stopMS() {
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

	var result map[string]interface{}
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
