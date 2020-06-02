package main

import (
	"github.com/alexflint/go-arg"
	"github.com/tkanos/gonfig"
)

//go get github.com/alexflint/go-arg
//go get github.com/go-resty/resty
//go get github.com/tkanos/gonfig

type configuration struct {
	IP       string
	Port     int
	Username string
	Password string
}

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
	config := configuration{}
	err := gonfig.GetConf("/etc/gologic.conf", &config)
	if err != nil {
		panic(err)
	}

	var args args
	arg.MustParse(&args)

	var admin = AdminServer{ipAdress: config.IP, port: config.Port, username: config.Username, password: config.Password}
	admin.init()

	switch args.Command {
	case "status":
		admin.statusAll(args.List)
	case "start":
		admin.start(args.List)
	case "stop":
		admin.stop(args.List)
	default:

	}
}
