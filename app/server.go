package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
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

func ParsePathForEcho(path string) (bool, string, string) {

	basePath, stringAfterEcho, echoPresent := strings.Cut(path, "echo/")

	fmt.Println("BasePath: "+basePath+"\nResponse String "+stringAfterEcho+"\nEcho Present ? ", echoPresent)

	return echoPresent, stringAfterEcho, basePath
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

	// httpResponse := "HTTP/1.1 200 OK\r\n\r\n"
	defaultResponse := "HTTP/1.1 404 Not Found\r\n\r\n"
	// echoPathResponse := "HTTP/1.1 200 OK\r\n\r\nContent-Type: text/plain\r\n\r\n"

	if method == "GET" {

		echoPresent, stringAfterEcho, _ := ParsePathForEcho(path)

		if echoPresent {

			contentLength := strconv.Itoa(len(stringAfterEcho))

			echoPathResponse := "HTTP/1.1 200 OK\r\n"
			echoPathResponse += "Content-type: text/plain\r\n"
			echoPathResponse += "Content-length: " + contentLength + "\r\n\r\n"
			echoPathResponse += stringAfterEcho + "\r\n\r\n"

			bytesSent, err := conn.Write([]byte(echoPathResponse))

			if err != nil {
				fmt.Println("Error sending response: ", err.Error())
				os.Exit(1)
			}
			fmt.Printf("Sent %d bytes to client (expected: %d)\n", bytesSent, len(echoPathResponse))
		} else {
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
