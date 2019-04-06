package main

import (
	"./conf"
	"./ssh"
	"flag"
	"fmt"
	"log"
	"sync"
)

var (
	group   *string = flag.String("g", "", "group name")
	cmd     *string = flag.String("c", "", "command")
	files   *string = flag.String("f", "", "source file")
	dst     *string = flag.String("d", "", "target  address")
	pwdfile *string = flag.String("p", ".auth.cfg", "Certification documents")
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
	flag.Parse()
	if *group == "" {
		flag.Usage()
		return
	}
	InfoList := conf.ReadConfig(*pwdfile, *group)

	wg := sync.WaitGroup{}
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
			continue
		}
		wg.Add(1)

		switch {
		case *cmd != "":
			go conn.Cmd(*cmd, &wg)
		case *files != "" && *dst != "":
			go conn.Scp(*files, *dst, &wg)
		case *cmd != "" && *files != "":
			flag.Usage()
			return
		default:
			flag.Usage()
			return
		}
	}
	defer wg.Wait()
}
