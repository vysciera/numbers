package transport

import (
	"fmt"
	"io"
	"net"
)

type UDP struct {
	conn   *net.UDPConn
	remote *net.UDPAddr
}

func NewUDP(address string) (*UDP, error) {
	remote, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, fmt.Errorf("resolve UDP address: %w", err)
	}

	// unconnected UDP socket
	// OS chooses the local address and ephemeral port.
	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return nil, fmt.Errorf("open UDP socket: %w", err)
	}

	return &UDP{
		conn:   conn,
		remote: remote,
	}, nil
}

func (u *UDP) Send(data []byte) error {
	n, err := u.conn.WriteToUDP(data, u.remote)
	if err != nil {
		return fmt.Errorf("UDP send: %w", err)
	}

	if n != len(data) {
		return io.ErrShortWrite
	}

	return nil
}

func (u *UDP) Close() error {
	return u.conn.Close()
}
