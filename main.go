package main

import (
	"./conf"
	"./ssh"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/urfave/cli"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ltime | log.Lshortfile)
}

var (
	grep    string
	cmd     string
	files   string
	dst     string
	pwdfile string
)

func main() {
	get_Args()
	InfoList := conf.ReadConfig(pwdfile, grep)

	wg := sync.WaitGroup{}
	wg.Add(len(InfoList))

	for _, info := range InfoList {
		conn := ssh.InfoSSH{
			User:      info.User,
			Password:  info.Password,
			PublicKey: info.PublicKey,
			Host:      info.Host,
			Port:      info.Port,
		}
		err := conn.Connect()
		if err != nil {
			fmt.Printf("\n \033[0;31m ==================== %v =======================  \033[0m\n", info.Host)
			fmt.Println(err)
			wg.Done()
			continue
		}

		switch {
		case cmd != "":
			go conn.Cmd(cmd, &wg)
		case files != "":
			go conn.Scp(files, dst, &wg)
		case cmd != "" && files != "":
			os.Exit(0)
		default:
			os.Exit(0)
		}
	}
	defer wg.Wait()
}

func get_Args() {
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
