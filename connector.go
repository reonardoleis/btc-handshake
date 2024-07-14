package main

import (
	"log"
	"net"
)

func connect(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println("error while opening connection with server", err)
		return nil, err
	}

	return conn, nil
}
