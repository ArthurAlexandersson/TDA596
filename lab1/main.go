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
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Handling connection")
}

// Main function.
func main() {

	var portNum string
	flag.StringVar(&portNum, "port", "8080", "Enter the port number which shoul be used for the server.")
	flag.Parse()

	fmt.Println("Starting http server with port number", portNum)
	_, input_err := fmt.Scan(&portNum)
	if input_err != nil {
		log.Fatal("Error reading from the command line.")
		return
	}

	log.Println("Started http server with port number", portNum)
}
