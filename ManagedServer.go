package main

import (
	"fmt"

	"github.com/go-resty/resty"
)

type ManagedServer struct {
	Name           string
	Status         string
	WeblogicHome   string
	MiddlewareHome string
	Cli            *resty.Client
}

func (ms *ManagedServer) statusMS() string {
	return ms.Status
}

func (ms *ManagedServer) startMS() {
	resp, err := ms.Cli.R().
		SetPathParams(map[string]string{
			"managedServerName": ms.Name,
		}).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Requested-By", "gologic").
		SetHeader("Prefer", "respond-async").
		SetBody(`{}`).
		Post("/management/weblogic/latest/domainRuntime/serverLifeCycleRuntimes/{managedServerName}/start")

	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

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
		SetBody(`{}`).
		Post("/management/weblogic/latest/domainRuntime/serverLifeCycleRuntimes/{managedServerName}/shutdown")

	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

}
