package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	// fork and make child process that constantly checks for http requests

	// Create a buffer to read the request line and headers
	reader := bufio.NewReader(conn)

	// Read the first line (Request-Line)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		// return error 400 Bad Request
		return
	}

	// Parse the request method, path, and protocol
	requestParts := strings.Fields(requestLine)
	if len(requestParts) < 3 {
		// return error 400 Bad Request
		return
	}

	method := requestParts[0]
	path := requestParts[1]

	// Check the method and handle GET or POST
	switch method {
	case "GET":
		handleGetRequest(conn, path)
	default:
		// For unsupported methods
		// Return error 400 Bad Request
	}
}

func handleGetRequest(conn net.Conn, path string) {

	// Check if the file name contains a dot (i.e., it should have an extension)
	copyPath := path
	if !strings.Contains(copyPath, ".") {
		// Handle the case where no dot is found (no file extension)
		response := &httpResponseWriter{conn: conn}
		response.WriteHeader(http.StatusBadRequest) // Send 400 Bad Request
		response.Write([]byte("400 Bad Request"))   // Provide an error message in the response body
		return
	}
	// Remove leading slash from the path
	fileName := strings.TrimPrefix(path, "/")
	fileType := filepath.Ext(fileName)[1:] // Extract file extension without the dot

	// Use a ResponseWriter to write back to the client
	response := &httpResponseWriter{conn: conn}
	GET(response, fileType, fileName)
}

// Main function.
func main() {

	var portNum string
	flag.StringVar(&portNum, "port", "8080", "Enter the port number which shoul be used for the server.")
	flag.Parse()
	fmt.Println("Starting server on port: ", portNum)
	OwnListenAndServe(portNum) // Start listening on the port number.
}

// Helper struct to implement http.ResponseWriter
type httpResponseWriter struct {
	conn        net.Conn
	headers     http.Header
	statusCode  int
	wroteHeader bool
}

// Initialize the header map in the `httpResponseWriter`
func (w *httpResponseWriter) Header() http.Header {
	if w.headers == nil {
		w.headers = make(http.Header)
	}
	return w.headers
}

// Modify `WriteHeader` to include a properly formatted status line and headers
func (w *httpResponseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.statusCode = statusCode
	w.wroteHeader = true

	// Write the status line
	fmt.Fprintf(w.conn, "HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))

	// Write headers
	for key, values := range w.Header() {
		for _, value := range values {
			fmt.Fprintf(w.conn, "%s: %s\r\n", key, value)
		}
	}
	fmt.Fprintf(w.conn, "\r\n") // End of headers
}

// Write the response body
func (w *httpResponseWriter) Write(data []byte) (int, error) {
	// Set Content-Length if not already set
	if w.Header().Get("Content-Length") == "" {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	}
	// Write headers if they haven’t been written
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.conn.Write(data)
}

func GET(w http.ResponseWriter, fileType string, fileName string) {
	// Map file type to MIME type
	contentType := ""

	switch fileType {
	case "html":
		contentType = "text/html"
	case "txt":
		contentType = "text/plain"
	case "gif":
		contentType = "image/gif"
	case "jpeg", "jpg":
		contentType = "image/jpeg"
	case "css":
		contentType = "text/css"
	default:
		http.Error(w, "Error 400: Unsupported file type", http.StatusBadRequest)
		return
	}

	// Construct full file path
	filePath := filepath.Join("files", fileName)
	// Check if file exists and is accessible
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Error 404: File not found", http.StatusNotFound)
		return
	}
	// Read the file content
	fileContent, _ := ioutil.ReadFile(filePath)

	// Write content-type and file content to the response
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(fileContent)
}
