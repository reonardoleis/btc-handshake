package handshake

import (
	"bytes"
	"log"
	"time"
)

func Handshake(nodeAddress string) {
	conn, err := connect(nodeAddress)
	if err != nil {
		log.Println("handshake: error while estabilishing connection with the node", err)
		return
	}

	defer conn.Close()

	versionMessage := NewVersionMessage("/rxonvrdo@challenge:1.0.0/", 0, time.Now())
	var encodedMsgBuff bytes.Buffer
	err = versionMessage.Encode(&encodedMsgBuff)
	if err != nil {
		log.Println("handshake: error while encoding version message", err)
		return
	}

	err = send(conn, "version", encodedMsgBuff.Bytes())
	if err != nil {
		log.Println("handshake: error while sending version message", err)
		return
	}
	log.Println("handshake: version message sent")

	msg, err := receive(conn)
	if err != nil {
		log.Println("handshake: error while reading message", err)
		return
	}
	log.Println("handshake: received", msg.command)

	if valid := msg.validChecksum(); !valid {
		log.Println("handshake: invalid checksum, aborting...")
		return
	}

	err = send(conn, "verack", nil)
	if err != nil {
		log.Println("handshake: error sending verack message", err)
		return
	}
	log.Println("handshake: verack message sent")

	msg, err = receive(conn)
	if err != nil {
		log.Println("handshake: error while reading message", err)
		return
	}
	log.Println("handshake: received", msg.command)

	if valid := msg.validChecksum(); !valid {
		log.Println("handshake: invalid checksum, aborting...")
		return
	}

	msg, err = receive(conn)
	if err != nil {
		log.Println("handshake: error while reading message", err)
		return
	}
	log.Println("handshake: received", msg.command)

	if valid := msg.validChecksum(); !valid {
		log.Println("handshake: invalid checksum, aborting...")
		return
	}

	log.Println("hanshake: handshake done, closing connection")
}
