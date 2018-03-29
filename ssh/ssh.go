package ssh

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

type InfoSSH struct {
	User     string
	Password string
	Host     string
	Port     int
	Csession *ssh.Session
	Fsession *sftp.Client
}

func (self *InfoSSH) Cmd(cmd string, wg *sync.WaitGroup) {
	defer self.Csession.Close()
	defer wg.Done()
	fmt.Printf("\n \033[1;32m ==================== %v ======================= \033[0m\n", self.Host)
	self.Csession.Stdout = os.Stdout
	self.Csession.Stderr = os.Stderr
	self.Csession.Run(cmd)
}

func (self *InfoSSH) Scp(src, dst string, wg *sync.WaitGroup) {
	var remoteFileName string

	defer self.Csession.Close()
	defer self.Fsession.Close()
	defer wg.Done()

	if strings.HasSuffix(dst, "/") == false {
		dst = dst + "/"
	}

	remoteFileName = path.Base(src)
	dstFile, err := self.Fsession.Create(path.Join(path.Dir(dst), remoteFileName))
	if err != nil {
		fmt.Printf("\n \033[0;31m ==================== %v ======================= \033[0m\n", self.Host)
		fmt.Println(err)
		return
	}
	defer dstFile.Close()

	fileByte, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Printf("\n \033[0;31m ==================== %v ======================= \033[0m\n", self.Host)
		fmt.Println(err)
		return
	}
	dstFile.Write(fileByte)
	fmt.Printf("\n \033[1;32m ==================== %v ======================= \033[0m\n", self.Host)
	fmt.Println("File sent successfully")
}

func (self *InfoSSH) Connect() error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshclient    *ssh.Client
		fsession     *sftp.Client
		csession     *ssh.Session
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(self.Password))

	clientConfig = &ssh.ClientConfig{
		User:    self.User,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr = fmt.Sprintf("%s:%d", self.Host, self.Port)
	if sshclient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}

	if csession, err = sshclient.NewSession(); err != nil {
		return err
	}

	if fsession, err = sftp.NewClient(sshclient); err != nil {
		return err
	}

	self.Csession = csession
	self.Fsession = fsession
	return nil
}
