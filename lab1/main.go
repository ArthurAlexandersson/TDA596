package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

// Port we listen on.

// Handler functions.
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage")
}

func Info(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Info page")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello page")
}

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
	http.HandleFunc("/", Home)
	http.HandleFunc("/info", Info)
	http.HandleFunc("/hello", Hello)
	// Write to connection
}

// ^^ Handler functions. ^^

// Main function.
func main() {

	log.Println("Starting our simple http server.")

	// Take in the port number from command line.
	var portNum string
	fmt.Println("Enter the port number")

	_, input_err := fmt.Scan(&portNum)

	if input_err != nil {
		log.Fatal("Error reading from the command line.")
		return
	}

	log.Println("Started on port", portNum)

	fmt.Println("To close connection CTRL+C :-)")

	OwnListenAndServe(portNum)

	// Registering our handler functions, and creating paths.

	// Spinning up the server.

}
