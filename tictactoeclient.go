package main

import (
	// "bufio"
	"fmt"
	"log"
	"net"
	// "os"
	"github.com/arjunkrishnababu96/tictactoe"
	"strings"
)

func main() {
	fmt.Println("I'm client!")
	conn, err := net.Dial("tcp", "127.0.0.1:7776")
	if err != nil {
		log.Fatalf(" main: ", err)
	}
	defer conn.Close()

	n, err := playTicTacToe(conn)
	if err != nil	{
		log.Fatalf(" main() n=%v: %v", n, err)
	}

	// input := bufio.NewScanner(os.Stdin)
	// input.Scan()
	//
	// n, err := sendMsg(conn, input.Text())
	// if err != nil {
	// 	log.Fatalf(" error while writing: ", err)
	// }
	// fmt.Printf(" %v bytes written\n", n)
	//
	// var msg string
	// msg, n, err = receiveMsg(conn)
	// if err != nil {
	// 	log.Fatalf(" error while receiving: ", err)
	// }
	// fmt.Println("Server says: ", msg)
}

func playTicTacToe(conn net.Conn) (int, error)  {
	const CLIENTSYMBOL = 'X'
	squares := []int{0,1,2,4,5,6,8,9,10}
	board := tictactoe.GetEmptyBoard()

	for {
		board, _ = tictactoe.MakeRandomMove(board, squares, CLIENTSYMBOL)

		n, err := conn.Write([]byte(board))
		if err != nil	{
			return n, fmt.Errorf("playTicTacToe error while writing %v", board)
		}

		responseBytes := make([]byte, 100)
		n, err = conn.Read(responseBytes)
		if strings.Contains(string(responseBytes), "END")	{
			break
		}
	}
	return 0, nil
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
