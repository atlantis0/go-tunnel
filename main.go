package main

import (
	"flag"
	"log"
	"net"
)

var (
	ipAddress       = flag.String("ipAddress", "127.0.0.1", "Please attach a valid IP address")
	port            = flag.String("port", "1212", "The port you wish to use")
	remoteIPAddress = flag.String("remoteIPAddress", "127.0.0.1", "Please attach a valid IP address")
	remotePort      = flag.String("remotePort", "1234", "The port you wish to use")
	publicKey       = flag.String("publicKey", "public.pem", "Please enter the path of your public key")
	privateKey      = flag.String("privateKey", "private.pem", "Please enter the path of your private key")
	createCert      = flag.Bool("createCert", false, "Create Public and Private PEM")
	server          = flag.Bool("server", false, "You are accepting TLS connections from other hosts")
	client          = flag.Bool("client", false, "You are tunneling connections to a server")
)

func main() {

	var listener net.Listener
	var secure bool
	flag.Parse()

	if *createCert {
		CreateEncryptionKeys()
		log.Println("TLS Certs Created")
		return
	}

	if *client == false && *server == false {
		flag.Usage()
		return
	}

	// this is the binding port ;)
	if addr := net.ParseIP(*ipAddress); addr == nil {
		log.Fatalln("Unable to parse IP. IP Provided was", *ipAddress)
	}
	service := *ipAddress + ":" + *port

	// remote ip address to connect to
	if raddr := net.ParseIP(*remoteIPAddress); raddr == nil {
		log.Fatalln("Unable to parse IP. IP Provided was", *remoteIPAddress)
	}
	RemoteIPandPort := *remoteIPAddress + ":" + *remotePort

	if *server {
		// server mode ;)
		log.Println("Server mode ...")
		listener = ServeTLSConnections(*publicKey, *privateKey, service)
		secure = false // We will be connecting to our service using TCP
	} else { // if we are not a server we must be a client
		listener = ServeTCPConnections(service)
		secure = true // We need to connect using TLS over the wire
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Println("Connection Accepted from", conn.RemoteAddr().String())
		// handle client
		go handleClient(conn, RemoteIPandPort, secure)
	}
}
