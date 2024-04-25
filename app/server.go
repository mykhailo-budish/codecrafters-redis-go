package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var readError error

	buffer := make([]byte, 1<<12)
	for readError != io.EOF {
		bytesRead, readError := conn.Read(buffer)
		if readError == io.EOF {
			break
		}
		if readError != nil {
			panic(readError)
		}

		command := ParseCommand(buffer[:bytesRead])
	}
}
