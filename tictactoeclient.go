package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("I'm client!")
	conn, err := net.Dial("tcp", "127.0.0.1:7776")
	if err != nil {
		log.Fatalf(" main: ", err)
	}
	defer conn.Close()

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	n, err := sendMsg(conn, input.Text())
	if err != nil {
		log.Fatalf(" error while writing: ", err)
	}
	fmt.Printf(" %v bytes written\n", n)

	var msg string
	msg, n, err = receiveMsg(conn)
	if err != nil {
		log.Fatalf(" error while receiving: ", err)
	}
	fmt.Println("Server says: ", msg)
}

func sendMsg(conn net.Conn, msg string) (n int, err error) {
	n, err = conn.Write([]byte(msg))
	if err != nil {
		return n, fmt.Errorf(" error while writing: ", err)
	}
	return n, nil
}

func receiveMsg(conn net.Conn) (msg string, n int, err error) {
	var msg_bytes = make([]byte, 100)
	n, err = conn.Read(msg_bytes)
	if err != nil {
		return msg, n, fmt.Errorf(" receiveMsg: %v", err)
	}

	return string(msg_bytes), n, err
}
