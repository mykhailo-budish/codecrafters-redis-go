package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	for i := 0; i < 5; i++ {
		conn, err := net.Dial("tcp", "0.0.0.0:6379")
		if err != nil {
			fmt.Println("Failed to bind to port 6379")
			os.Exit(1)
		}

		conn.Write([]byte("*1\r\n$4\r\nping\r\n"))

		buf := make([]byte, 1024)
		bytesRead, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(buf[:bytesRead]))
		conn.Close()
	}
}
