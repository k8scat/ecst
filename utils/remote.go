package utils

import (
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

const (
	DefaultSSHPort     = "22"
	DefaultSSHUser     = "root"
	DefaultPermissions = "0777"
	DefaultDir         = "/tmp"
)

func NewSSHClient(user, password, ip, port string) (client *ssh.Client, err error) {
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * time.Duration(2),
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
	}
	addr := fmt.Sprintf("%s:%s", ip, port)
	client, err = ssh.Dial("tcp", addr, config)
	return
}

func NewScpClient(user, password, ip, port string) (client scp.Client, err error) {
	config, _ := auth.PasswordKey(user, password, ssh.InsecureIgnoreHostKey())
	addr := fmt.Sprintf("%s:%s", ip, port)
	return scp.NewClient(addr, &config), nil
}

func RunCommandRemote(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	return session.Run(command)
}

func Scp(client scp.Client, localFile, remotePath, permissions string) error {
	err := client.Connect()
	if err != nil {
		return err
	}
	f, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return client.CopyFile(f, remotePath, permissions)
}

func RunScriptRemote(user, password, ip, port, localFile, remotePath, permissions string) error {
	scpClient, err := NewScpClient(user, password, ip, port)
	if err != nil {
		return err
	}
	defer scpClient.Close()
	err = Scp(scpClient, localFile, remotePath, permissions)
	if err != nil {
		return err
	}
	sshClient, err := NewSSHClient(user, password, ip, port)
	if err != nil {
		return err
	}
	defer sshClient.Close()
	command := fmt.Sprintf("/bin/sh %s", remotePath)
	return RunCommandRemote(sshClient, command)
}
