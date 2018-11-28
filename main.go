package main

import (
	"./conf"
	"./ssh"
	"fmt"
	"log"
	"runtime"
	"sync"
)

func init() {
    GetArgs()
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
	InfoList := conf.ReadConfig(pwdfile, grep)

	wg := sync.WaitGroup{}
	wg.Add(len(InfoList))

	for _, info := range InfoList {
		conn := &ssh.InfoSSH{
			User:      info.User,
			Password:  info.Password,
			PublicKey: info.PublicKey,
			Host:      info.Host,
			Port:      info.Port,
		}
        if err := conn.Connect(); err != nil {
			fmt.Printf("\n \033[0;31m ==================== %v =======================  \033[0m\n", info.Host)
			fmt.Println(err)
			wg.Done()
			continue
		}

		switch {
		case cmd != "":
			go conn.Cmd(cmd, &wg)
		case files != "" && dst != "":
			go conn.Scp(files, dst, &wg)
		case cmd != "" && files != "":
			log.Fatal("Parameter error")
		default:
			log.Fatal("Parameter error")
		}
	}
	defer wg.Wait()
}
