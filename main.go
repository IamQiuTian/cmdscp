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
	grep  string
	cmd   string
	files string
	dst   string
)

func main() {
	get_Args()

	InfoList := conf.ReadConfig(grep)

	wg := sync.WaitGroup{}
	wg.Add(len(InfoList))
 
	for _, info := range InfoList {
		conn := ssh.InfoSSH{
			User:     info.User,
			Password: info.Password,
			Host:     info.Host,
			Port:     info.Port,
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
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "grep,g",
			Value:       "",
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
	}
	app.Action = func(c *cli.Context) error {
		if c.String("grep") == "" {
			cli.ShowSubcommandHelp(c)
			os.Exit(0)
		}
		return nil
	}
	app.Run(os.Args)
}
