package handshake

import (
	"bytes"
	"encoding/binary"
	"net"
)

type message struct {
	command  string
	payload  []byte
	checksum [4]byte
}

// reads 24 bytes from the connection, which represents the header of the message
// reads 4 bytes from 16 to 20 from the header, to get the payload length, and
// then reads that amount of bytes from the connection, in order to get the payload.
// gets the command by trimming the header between 4 and 16
// gets the checksum by copying header between 20 and 24
func receive(conn net.Conn) (message, error) {
	header := make([]byte, 24)
	_, err := conn.Read(header)
	if err != nil {
		return message{}, err
	}

	payloadLength := binary.LittleEndian.Uint32(header[16:20])
	payload := make([]byte, payloadLength)
	_, err = conn.Read(payload)
	if err != nil {
		return message{}, err
	}

	command := string(bytes.Trim(header[4:16], "\x00"))

	var checksum [4]byte
	copy(checksum[:], header[20:24])

	return message{command: command, payload: payload, checksum: checksum}, nil
}

func (m message) validChecksum() bool {
	return bytes.Equal(m.checksum[:], doubleSha256(m.payload)[:4])
}
