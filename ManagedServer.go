package main

import "github.com/go-resty/resty"

type ManagedServer struct {
	status         string
	weblogicHome   string
	middlewareHome string
	cli            *resty.Client
}
