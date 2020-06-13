package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/spf13/cobra"
	"github.com/tkanos/gonfig"
)

//go get github.com/alexflint/go-arg
//go get github.com/go-resty/resty
//go get github.com/tkanos/gonfig
//"github.com/spf13/cobra"

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
	var configPath string = "/etc/gologic.conf"
	config := configuration{}
	err := gonfig.GetConf(configPath, &config)
	if err != nil {
		panic(err)
	}

	var args args
	arg.MustParse(&args)

	var admin = AdminServer{ipAdress: config.IP, port: config.Port, username: config.Username, password: config.Password}
	admin.init()

	switch args.Command {
	case "start":
		admin.start(args.List)
	case "stop":
		admin.stop(args.List)
	case "info":
		admin.printInfo()
		fmt.Printf("%-40s %s\n", "Configfile", configPath)
	case "add":
		admin.createManagedServer(args.List[0], args.List[1], args.List[2])
	case "deployment":
		admin.printDeployments()
	default:
		fmt.Println("Usage guide")
	}

	var cmdStatus = &cobra.Command{
		Use:   "status [string to status]",
		Short: "status anything to the screen",
		Long: `status is for printing anything back to the screen.
  For many years people have printed back to the screen.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(len(args))
			admin.printStatus(args)
		},
	}

	var rootCmd = &cobra.Command{Use: "gologic"}
	rootCmd.AddCommand(cmdStatus)
	rootCmd.Execute()
}
