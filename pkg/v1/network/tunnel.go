package network

import (
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

type TunnelPayload struct {
	SSHHost      string // e.g., "tunnel.example.com:22"
	SSHUser      string // e.g., "ubuntu"
	AuthType     string // "password" or "key"
	Secret       string // Password or path to private key (or raw private key string)
	LocalBind    string // e.g., "127.0.0.1:5432" (where your local app connects)
	RemoteTarget string // e.g., "internal-db.local:5432" (the destination behind the firewall)
}

func CreateTunnel(payload TunnelPayload) error {
	var authMethod ssh.AuthMethod
	switch payload.AuthType {
	case "password":
		authMethod = ssh.Password(payload.Secret)
	case "key":
		keyBytes, err := os.ReadFile(payload.Secret)
		if err != nil {
			return fmt.Errorf("failed to read private key file: %w", err)
		}
		signer, err := ssh.ParsePrivateKey(keyBytes)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethod = ssh.PublicKeys(signer)
	default:
		return fmt.Errorf("unsupported auth type: %s", payload.AuthType)
	}

	config := &ssh.ClientConfig{
		User:            payload.SSHUser,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use strict verification in production
	}

	// 1. Start listening on the local interface (container localhost or host loopback)
	listener, err := net.Listen("tcp", payload.LocalBind)
	if err != nil {
		return fmt.Errorf("failed to bind local port %s: %w", payload.LocalBind, err)
	}

	// 2. Handle incoming connections in a background listener loop
	go func() {
		defer listener.Close()
		for {
			localConn, err := listener.Accept()
			if err != nil {
				break
			}

			// Handle each connection concurrently
			go func(clientConn net.Conn) {
				defer clientConn.Close()

				// Dial to the SSH server
				sshConn, err := ssh.Dial("tcp", payload.SSHHost, config)
				if err != nil {
					return
				}
				defer sshConn.Close()

				// Open a channel through the SSH tunnel to the final remote destination
				remoteConn, err := sshConn.Dial("tcp", payload.RemoteTarget)
				if err != nil {
					return
				}
				defer remoteConn.Close()

				// Bidirectional data copy stream container <-> tunnel <-> destination
				closed := make(chan struct{}, 2)
				go func() {
					io.Copy(remoteConn, clientConn)
					closed <- struct{}{}
				}()
				go func() {
					io.Copy(clientConn, remoteConn)
					closed <- struct{}{}
				}()
				<-closed
			}(localConn)
		}
	}()

	return nil
}
