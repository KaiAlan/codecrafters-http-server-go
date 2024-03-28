package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func ParseRequest(request string) (string, string, string) {

	requestFirstLine := strings.Split(request, "\r\n")[0]
	requestParams := strings.Split(requestFirstLine, " ")

	method := requestParams[0]
	path := requestParams[1]
	version := requestParams[2]

	return method, path, version
}

func main() {

	fmt.Println("Logging...")

	//binding to port
	listner, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := listner.Accept()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	fmt.Println("Client connected")

	//reading incomig requests
	readBuffer := make([]byte, 2048)
	bytesReceived, err := conn.Read(readBuffer)

	if err != nil {
		fmt.Printf("Error reading request: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Read %d bytes from client\n", bytesReceived)

	request := string(readBuffer[:bytesReceived])

	//Paresing Request and Responsing with proper response

	method, path, _ := ParseRequest(request)

	httpResponse := "HTTP/1.1 200 OK\r\n\r\n"
	defaultResponse := "HTTP/1.1 404 Not Found\r\n\r\n"

	if method == "GET" {

		switch path {
		case "/":
			bytesSent, err := conn.Write([]byte(httpResponse))

			if err != nil {
				fmt.Println("Error sending response: ", err.Error())
				os.Exit(1)
			}
			fmt.Printf("Sent %d bytes to client (expected: %d)\n", bytesSent, len(httpResponse))
		default:
			_, err := conn.Write([]byte(defaultResponse))

			if err != nil {
				fmt.Println("Error sending response: ", err.Error())
				os.Exit(1)
			}
		}
	} else {
		fmt.Println("Not a GET request")
	}

}
