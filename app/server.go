package main

import (
	"flag"
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

func setResponse(statusCode string, statusMessage string, contentType string, contentLength string, responseBody string) string {
	response := "HTTP/1.1 " + statusCode + " " + statusMessage + "\r\n"
	response += "Content-type: " + contentType + "\r\n"
	response += "Content-length: " + contentLength + "\r\n\r\n"
	response += responseBody + "\r\n"

	return response
}

func sendResponse(response string, conn net.Conn) {
	bytesSent, err := conn.Write([]byte(response))

	if err != nil {
		fmt.Println("Error sending response: ", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Sent %d bytes to client (expected: %d)\n", bytesSent, len(response))
}

func main() {

	fmt.Println("Logging...")

	filePath := flag.String("directory", "", "file path")
	flag.Parse()

	//binding to port
	listner, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer listner.Close()

	for {

		conn, err := listner.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Client connected: ", conn.RemoteAddr())

		go handleConnection(conn, *filePath)
	}

}

func handleConnection(conn net.Conn, filePath string) {

	defer conn.Close()

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

		if path == "/" {
			_, err := conn.Write([]byte(httpResponse))

			if err != nil {
				fmt.Println("Error sending response: ", err.Error())
				os.Exit(1)
			}
		} else if strings.HasPrefix(path, "/echo") {

			stringAfterEcho := strings.TrimLeft(path, "/echo/")
			contentLength := strconv.Itoa(len(stringAfterEcho))

			response := setResponse("200", "OK", "text/plain", contentLength, stringAfterEcho)

			sendResponse(response, conn)
		} else if path == "/user-agent" {

			requestHeader := strings.Split(request, "\r\n")[2]

			if strings.HasPrefix(requestHeader, "User-Agent") {

				requestBody := requestHeader[12:]
				response := setResponse("200", "OK", "text/plain", strconv.Itoa(len(requestBody)), requestBody)
				sendResponse(response, conn)
			}
		} else if strings.HasPrefix(path, "/files") {

			filename := strings.Split(path, "/files/")[1]
			fullFilePath := filePath + "/" + filename

			fileContent, err := os.ReadFile(fullFilePath)

			if err != nil {
				_, err := conn.Write([]byte(defaultResponse))

				if err != nil {
					fmt.Println("Error sending response: ", err.Error())
					os.Exit(1)
				}
				os.Exit(1)
			}

			response := setResponse("200", "OK", "application/octet-stream", strconv.Itoa(len(fileContent)), string(fileContent))
			sendResponse(response, conn)

		} else {
			_, err := conn.Write([]byte(defaultResponse))

			if err != nil {
				fmt.Println("Error sending response: ", err.Error())
				os.Exit(1)
			}
		}
	} else if method == "POST" {

		if strings.HasPrefix(path, "/files/") {
			filename := strings.Split(path, "/files/")[1]
			fullFilePath := filePath + "/" + filename
			requestBody := strings.Split(request, "\r\n\r\n")[1]

			err := os.WriteFile(fullFilePath, []byte(requestBody), 0666)

			if err == nil {

				fmt.Println("Writing to file: " + fullFilePath + " (Total content length: " + strconv.Itoa(len(requestBody)) + ")")

				conn.Write([]byte("HTTP/1.1 201 OK\r\n\r\n"))
			} else {
				fmt.Println("Server Error: ", err.Error())
				os.Exit(1)
			}
		}
	} else {
		fmt.Println("Not a valid request")
	}

}
