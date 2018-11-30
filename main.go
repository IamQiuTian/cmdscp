package main

import (
	"./conf"
	"./ssh"
	"fmt"
	"log"
    "flag"
	"sync"
)

var (
    Group   *string = flag.String("g", "", "group name")
    Cmd     *string = flag.String("c", "", "command")
    Files   *string = flag.String("f", "", "source file")
    Dst     *string = flag.String("d", "", "target  address")
    Pwdfile *string = flag.String("p", ".info.json", "Certification documents")
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
    flag.Parse()
    if *Group == "" {
        flag.Usage()
        return
    }
	InfoList := conf.ReadConfig(*Pwdfile, *Group)

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
		case *Cmd != "":
			go conn.Cmd(*Cmd, &wg)
		case *Files != "" && *Dst != "":
			go conn.Scp(*Files, *Dst, &wg)
		case *Cmd != "" && *Files != "":
			flag.Usage()
            return
		default:
			flag.Usage()
            return
		}
	}
	defer wg.Wait()
}
