package handshake

import (
	"encoding/binary"
	"net"
)

// sends the header and the payload
// the header is 24 bytes long, consisting of:
// - 4 bytes to identify the network (mainnet, testnet, etc)
// - 12 bytes to store the command (version, verack, ping, etc), filling with 0s if it is not used fully
// - 4 bytes to store the length of the payload
// - 4 bytes to store the doubleSha256 value of the first 4 bytes of the payload
func send(conn net.Conn, command string, payload []byte) error {
	header := make([]byte, 24)

	binary.LittleEndian.PutUint32(header[0:4], uint32(Mainnet))

	copy(header[4:16], append([]byte(command), make([]byte, 12-len(command))...))

	binary.LittleEndian.PutUint32(header[16:20], uint32(len(payload)))

	copy(header[20:24], doubleSha256(payload)[:4])

	_, err := conn.Write(header)
	if err != nil {
		return err
	}

	_, err = conn.Write(payload)
	if err != nil {
		return err
	}

	return nil
}
