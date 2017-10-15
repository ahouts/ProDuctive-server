package tunnel

import (
	"fmt"
	"io"
	"net"

	"golang.org/x/crypto/ssh"
)

type SSHTunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

func (tunnel *SSHTunnel) Start(ready chan (bool)) error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()
	ready <- true
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
		return
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}
