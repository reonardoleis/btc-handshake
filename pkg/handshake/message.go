package handshake

import (
	"encoding/binary"
	"io"
	"log"
	"math"
	"net"
	"time"
)

type NetAddress struct {
	Timestamp time.Time
	Services  uint64
	IP        net.IP
	Port      uint16
}

type VersionMessage struct {
	ProtocolVersion ProtocolVersion
	Services        ServiceFlag
	Timestamp       time.Time
	AddrYou         NetAddress
	AddrMe          NetAddress
	Nonce           uint64
	UserAgent       string
	LastBlock       int32
	DisableRelayTx  bool
}

func NewVersionMessage(userAgent string, lastBlock int32, timestamp time.Time) *VersionMessage {
	return &VersionMessage{
		ProtocolVersion: 70016,
		UserAgent:       userAgent,
		LastBlock:       lastBlock,
		Services:        FullNodeFlag,
		Timestamp:       timestamp,
	}
}

func (m *VersionMessage) writeNetAdresses(w io.Writer) error {
	buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, m.AddrYou.Services)
	if _, err := w.Write(buff); err != nil {
		return err
	}

	var ip [16]byte
	if m.AddrYou.IP != nil {
		copy(ip[:], m.AddrYou.IP.To16())
	}
	if _, err := w.Write(ip[:]); err != nil {
		return err
	}

	binary.BigEndian.PutUint16(buff[:2], m.AddrYou.Port)
	if _, err := w.Write(buff[:2]); err != nil {
		return err
	}

	buff = make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, m.AddrMe.Services)
	if _, err := w.Write(buff); err != nil {
		return err
	}

	ip = [16]byte{}
	if m.AddrMe.IP != nil {
		copy(ip[:], m.AddrMe.IP.To16())
	}
	if _, err := w.Write(ip[:]); err != nil {
		return err
	}

	binary.BigEndian.PutUint16(buff[:2], m.AddrMe.Port)
	if _, err := w.Write(buff[:2]); err != nil {
		return err
	}

	return nil
}

func (m *VersionMessage) writeUserAgent(w io.Writer) error {
	buff := make([]byte, 5) // hold up to a uint32 with 1 extra byte to mark if it is a uint64
	strlen := uint64(len(m.UserAgent))

	switch {
	case strlen < 0xfd:
		buff[0] = uint8(strlen)
		if _, err := w.Write(buff[:1]); err != nil {
			return err
		}

	case strlen <= math.MaxUint16:
		buff[0] = 0xfd
		binary.LittleEndian.PutUint16(buff[1:3], uint16(strlen))
		if _, err := w.Write(buff[:3]); err != nil {
			return err
		}

	case strlen <= math.MaxUint32:
		buff[0] = 0xfe
		binary.LittleEndian.PutUint32(buff[1:5], uint32(strlen))
		if _, err := w.Write(buff[:5]); err != nil {
			return err
		}

	default:
		buff[0] = 0xff
		if _, err := w.Write(buff[:1]); err != nil {
			return err
		}

		binary.LittleEndian.PutUint64(buff, strlen)
		if _, err := w.Write(buff); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte(m.UserAgent)); err != nil {
		return err
	}

	return nil
}

func (m *VersionMessage) Encode(w io.Writer) error {
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, uint32(m.ProtocolVersion))
	if _, err := w.Write(buff); err != nil {
		log.Println("handshake: error while encoding protocol version", err)
		return err
	}

	buff = make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, m.Services)
	if _, err := w.Write(buff); err != nil {
		log.Println("handshake: error while encoding services flag", err)
		return err
	}

	buff = make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, uint64(m.Timestamp.Unix()))
	if _, err := w.Write(buff); err != nil {
		log.Println("handshake: error while encoding timestamp", err)
		return err
	}

	if err := m.writeNetAdresses(w); err != nil {
		log.Println("handshake: error while encoding net addresses", err)
		return err
	}

	buff = make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, m.Nonce)
	if _, err := w.Write(buff); err != nil {
		log.Println("handshake: error while encoding nonce", err)
		return err
	}

	if err := m.writeUserAgent(w); err != nil {
		log.Println("handshake: error while encoding user agent", err)
	}

	buff = make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, uint32(m.LastBlock))
	if _, err := w.Write(buff); err != nil {
		log.Println("handshake: error while encoding last block", err)
		return err
	}

	return nil
}
