package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func OwnListenAndServe(port string) {
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Connection accepted")
		}
		go handleConnection(conn) // Handle the connection.

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Handling connection")
	// TODO Handle the connection.
}

// Main function.
func main() {

	var portNum string
	flag.StringVar(&portNum, "port", "8080", "Enter the port number which shoul be used for the server.")
	flag.Parse()
	fmt.Println("Starting server on port: ", portNum)
	OwnListenAndServe(portNum) // Start listening on the port number.
}
