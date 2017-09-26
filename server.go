package main

import (
	"crypto/rand"
	"crypto/tls"
	"log"
	"net"
	"time"
)

// ServeTCPConnections -
func ServeTCPConnections(service string) net.Listener {
	listener, err := net.Listen("tcp", service)
	checkError(err, true)
	log.Println("Listening for TCP connections on", service)
	return listener
}

// ServeTLSConnections -> server TLS Connections
// its now accepting connections at the given service point
// e.g. at localhost:8087
func ServeTLSConnections(publicKey, privateKey, service string) net.Listener {
	cert, err := tls.LoadX509KeyPair(publicKey, privateKey)
	checkError(err, true)

	config := tls.Config{Certificates: []tls.Certificate{cert}}

	now := time.Now()
	config.Time = func() time.Time { return now }
	config.Rand = rand.Reader

	//	listener, err := net.Listen("tcp", service)
	listener, err := tls.Listen("tcp", service, &config)
	checkError(err, true)
	log.Println("Listening for TLS connections on", service)
	return listener
}

// GetConnection - Get Connection - to destination
// this server should be reach the destination (in the same VPC or sth)
func GetConnection(remoteServerIPandPort string, secure bool) net.Conn {

	if !secure {
		serverConn, err := net.Dial("tcp", remoteServerIPandPort)
		checkError(err, true)
		// return connection
		return serverConn
	}

	cert, err := tls.LoadX509KeyPair("public.pem", "private.pem")
	checkError(err, true)

	config := tls.Config{Certificates: []tls.Certificate{cert}, Time: func() time.Time { return time.Now() },
		Rand: rand.Reader, InsecureSkipVerify: true}

	serverConn, err := tls.Dial("tcp", remoteServerIPandPort, &config)
	checkError(err, true)
	return serverConn

}

// handleClient - connecting clients are served here
// three parameters
// conn - conn accepted from the connecting client
// remoteServerIPandPort - remote application to reach
func handleClient(conn net.Conn, remoteServerIPandPort string, secure bool) {
	defer conn.Close()

	serverConn := GetConnection(remoteServerIPandPort, secure)
	defer serverConn.Close()

	var inBuf [10000]byte
	var outBuf [10000]byte

	for {
		// read from client
		n, err := conn.Read(inBuf[0:])
		if err != nil {
			log.Println("Error occured while reading incoming client!")
			log.Println(err)
			return
		}
		log.Println("read from client")

		// relay to the remote server
		_, err2 := serverConn.Write(inBuf[0:n])
		if err2 != nil {
			log.Println("Error while writting to the remote server!")
			return
		}
		log.Println("sent to remote server")

		// read response from remote server
		o, err3 := serverConn.Read(outBuf[0:])
		if err3 != nil {
			log.Println(err3)
			return
		}
		log.Println("read from remote server")

		// relay the response back to the client
		_, err4 := conn.Write(outBuf[0:o])
		if err4 != nil {
			return
		}
		log.Println("remote response to client")
	}
}
