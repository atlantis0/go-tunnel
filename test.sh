#!/bin/sh
rm GoTunnel # Cleanup
go build || exit
BIN_NAME="$(basename `pwd` )"
./$BIN_NAME --client &
./$BIN_NAME --server --port=1234 --remote-port=9000 &
echo "Starting Netcat on port 9000. Open another console and type in nc 127.0.0.1 1212"
nc -l 9000 
rm $BIN_NAME


# generate server
# var (
# 	ipAddress       = flag.String("ipAddress", "127.0.0.1", "Please attach a valid IP address")
# 	port            = flag.String("port", "1212", "The port you wish to use")
# 	remoteIPAddress = flag.String("remoteIPAddress", "127.0.0.1", "Please attach a valid IP address")
# 	remotePort      = flag.String("remotePort", "1234", "The port you wish to use")
# 	publicKey       = flag.String("publicKey", "public.pem", "Please enter the path of your public key")
# 	privateKey      = flag.String("privateKey", "private.pem", "Please enter the path of your private key")
# 	createCert      = flag.Bool("createCert", false, "Create Public and Private PEM")
# 	server          = flag.Bool("server", false, "You are accepting TLS connections from other hosts")
# 	client          = flag.Bool("client", false, "You are tunneling connections to a server")
# )

# to create certificates
# ./go-tunnel --ipAddress=127.0.0.1 --port=1212 --remoteIPAddress=54.202.175.229 --remotePort=3128 --createCert=True --server=true --client=false

# create a p12 file
# openssl pkcs12 -export -inkey private.key -in public.pem -name tunnel -out final_result.pfx
openssl pkcs12 -export -out local_certificate.p12 -inkey private.pem -in public.pem

# for running the server
# ./go-tunnel --ipAddress=127.0.0.1 --port=1212 --remoteIPAddress=54.202.175.229 --remotePort=3128 --createCert=False --server=True --client=False
