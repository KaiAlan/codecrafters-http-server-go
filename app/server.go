package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()

	defer conn.Close()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Client connected")
	readBuffer := make([]byte, 2048)
	bytesReceived, err := conn.Read(readBuffer)
	if err != nil {
		fmt.Printf("Error reading request: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Read %d bytes from client\n", bytesReceived)

	httpResponse := "HTTP/1.1 200 OK\r\n\r\n"
	bytesSent, err := conn.Write([]byte(httpResponse))
	if err != nil {
		fmt.Println("Error sending response: ", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Sent %d bytes to client (expected: %d)\n", bytesSent, len(httpResponse))
}
