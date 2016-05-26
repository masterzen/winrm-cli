package main

import (
	"fmt"
	"os"

	"github.com/masterzen/winrm-cli/command"
	"github.com/urfave/cli"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{
	{
		Name:   "gencert",
		Usage:  "Generate x509 client certificate to use with secure connections",
		Action: command.CmdGencert,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "size, s",
				Value: 2048,
				Usage: "size of the generated private key (possible values: 512, 1024, 2048, 4096)",
			},
		},
	},
	{
		Name:        "exec",
		Usage:       "Remotely execute a command",
		Description: "Argument is the command to remotely execute.",
		Action:      command.CmdExec,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "hostname, H",
				Value: "localhost",
				Usage: "WinRM host",
			},
			cli.StringFlag{
				Name:   "username, u",
				Value:  "vagrant",
				Usage:  "winrm admin username",
				EnvVar: "WINRM_USER",
			},
			cli.StringFlag{
				Name:   "password, p",
				Value:  "vagrant",
				Usage:  "winrm admin password",
				EnvVar: "WINRM_PASSWORD",
			},
			cli.IntFlag{
				Name:  "port, P",
				Value: 5985,
				Usage: "winrm port",
			},
			cli.BoolFlag{
				Name:  "https",
				Usage: "use https",
			},
			cli.BoolFlag{
				Name:  "insecure, i",
				Usage: "skip SSL validation",
			},
			cli.StringFlag{
				Name:  "cacert",
				Value: "",
				Usage: "Use CA certificates from `FILE`",
			},
			cli.StringFlag{
				Name:  "timeout",
				Value: "0s",
				Usage: "connection timeout",
			},
		},
	},
	{
		Name:   "version",
		Usage:  "Displays winrm version",
		Action: command.CmdVersion,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
