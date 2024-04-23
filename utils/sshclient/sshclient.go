package sshclient

import "golang.org/x/crypto/ssh"

func New(ip, username, password, port string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Ignore host key verification
	}

	return ssh.Dial("tcp", ip+":"+port, config)
}
