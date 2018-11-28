package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func GetArgs() {
	app := cli.NewApp()
	app.Name = "cmdscp"
	app.Version = "v0.0.1"
	app.Usage = "shell and send file"
	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "grep,g",
			Value:       "false",
			Usage:       "Input ip grep",
			Destination: &grep,
		},
		cli.StringFlag{
			Name:        "cmd,c",
			Value:       "",
			Usage:       "Input command",
			Destination: &cmd,
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
		cli.StringFlag{
			Name:        "pwdfile,p",
			Value:       ".info.json",
			Usage:       "Input passwd file path",
			Destination: &pwdfile,
		},
	}
	app.Action = func(c *cli.Context) {
		if c.String("grep") == "false" {
			cli.ShowAppHelp(c)
			os.Exit(0)
		}

		_, err := os.Stat(c.String("pwdfile"))
		if os.IsNotExist(err) {
			fmt.Println("Password file is none!")
			os.Exit(0)
		}
	}

	app.Run(os.Args)
}
