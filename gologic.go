package main

import (
	"github.com/alexflint/go-arg"
)

//go get github.com/alexflint/go-arg
//go get github.com/go-resty/resty

type args struct {
	Command string   `arg:"positional" help:"status, start, stop"`
	List    []string `arg:"positional" help:"empty or a list of Managed Server"`
}

func (args) Version() string {
	return "gologic 0.0.1"
}

func (args) Description() string {
	return "Gologic is a client for the Weblogics RESTful Management Services"
}

func main() {
	var args args
	arg.MustParse(&args)

	var admin = AdminServer{ipAdress: "127.0.0.1", port: 7001, username: "weblogic", password: "password123"}
	admin.init()

	switch args.Command {
	case "status":
		admin.status(args.List)
	case "start":
		admin.start(args.List)
	case "stop":
		admin.stop(args.List)
	default:

	}
}
