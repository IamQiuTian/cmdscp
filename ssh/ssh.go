package ssh

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type InfoSSH struct {
	User     string
	Password string
	Host     string
	Port     int
	Csession *ssh.Session
	Fsession *sftp.Client
}

func (self *InfoSSH) Cmd(cmd string, vv bool) {
	defer self.Csession.Close()
	if vv == true {
		self.Csession.Stdout = os.Stdout
		self.Csession.Stderr = os.Stderr
	}
	self.Csession.Run(cmd)
}

func (self *InfoSSH) Scp(src, dst string) {
	defer self.Csession.Close()
	defer self.Fsession.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	remoteFileName := path.Base(src)

	dstFile, err := self.Fsession.Create(path.Join(dst, remoteFileName))
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	rSrcFile := bufio.NewReader(srcFile)
	rSrcFile.Peek(rSrcFile.Buffered())

	var bufLine []byte
	for {
		buf, err := rSrcFile.ReadByte()
		if err == io.EOF {
			break
		}
		bufLine = append(bufLine, buf)
	}
	dstFile.Write(bufLine)
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
