package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	HttpServer("127.0.0.1", 8080)
}

func HttpServer(ip string, port int) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("Error listening: %+v\n", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Start listening...")

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting: %+v\n", err.Error())
			os.Exit(1)
		}
		// Handle connections
		go handleRequest(conn)
	}
}

func readFullData(conn net.Conn) ([]byte, error) {
	data := []byte{}
	notFinishRead := true
	for notFinishRead {
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				notFinishRead = false
				break
			}
			return nil, fmt.Errorf("error occur while reading, err: %+v", err.Error())
		}
		data = append(data, buffer[:n]...)
		if n < 4096 {
			notFinishRead = false
			break
		}
	}
	return data, nil
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recover: %+v\n", r)
		}
	}()
	data, err := readFullData(conn)

	if err != nil {
		fmt.Printf("readFullData err: %+v\n", err)
		return
	}

	request := NewRequest()
	err = request.Parse(data)
	if err != nil {
		fmt.Printf("parseRequest err: %+v\n", err)
		return
	}

	// fmt.Printf("request: %+v\n", request)

	fmt.Println("================================================================")
	// Send a response back to person contacting us.
	// r := []byte("HTTP/1.1 200 OK\r\nConnection: close\r\nContent-Type: text/html\r\nContent-Length: 19\r\n\r\n<h1>Hola Mundo</h1>")
	r := []byte(`HTTP/1.1 200 OK
Connection: close
Content-Type: text/json
Content-Length: 0

`)
	conn.Write(r)
}
