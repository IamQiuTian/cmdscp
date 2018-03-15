package main

import (
	"./conf"
	"./ssh"
	"log"
	"os"

	"github.com/urfave/cli"
)

var (
	target string
	cmd    string
	show   string
	files  string
	dst    string
)

func main() {
	get_Args()
	InfoMap := conf.FindInfo(target)

	conn := ssh.InfoSSH{
		User:     InfoMap["username"].(string),
		Password: InfoMap["password"].(string),
		Host:     InfoMap["ip"].(string),
		Port:     InfoMap["port"].(int),
	}
	err := conn.Connect()
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case cmd != "" && show == "false":
		conn.Cmd(cmd, false)
	case cmd != "" && show != "false":
		conn.Cmd(cmd, true)
	case files != "":
		conn.Scp(files, dst)
	case cmd != "" && files != "":
		os.Exit(0)
	default:
		os.Exit(0)
	}
}

func get_Args() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "target,t",
			Value:       "",
			Usage:       "Input ip",
			Destination: &target,
		},
		cli.StringFlag{
			Name:        "cmd,c",
			Value:       "",
			Usage:       "Input command",
			Destination: &cmd,
		},
		cli.StringFlag{
			Name:        "info,show",
			Value:       "false",
			Usage:       "Print command result",
			Destination: &show,
		},
		cli.StringFlag{
			Name:        "file,f",
			Value:       "",
			Usage:       "Input source file path",
			Destination: &files,
		},
		cli.StringFlag{
			Name:        "dst,d",
			Value:       "",
			Usage:       "Input target file path",
			Destination: &dst,
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.String("target") == "" {
			cli.ShowSubcommandHelp(c)
			os.Exit(0)
		}
		return nil
	}
	app.Run(os.Args)
}
