package main

import (
	"fmt"
	"github.com/arjunkrishnababu96/tictactoe"
	"log"
	"net"
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
	if err != nil {
		log.Fatalf(" main() n=%v: %v", n, err)
	}
}

func playTicTacToe(conn net.Conn) (int, error) {
	const (
		CLIENTWON = "client won"
		SERVERWON = "server won"

		CLIENTSYMBOL = 'X'
		SERVERSYMBOL = 'O'
	)

	squares := []int{0, 1, 2, 4, 5, 6, 8, 9, 10}

	var rboard string
	var sboard string = tictactoe.GetEmptyBoard()

	var n int
	var err error

	// make first move before the infinite loop starts
	sboard, _ = tictactoe.MakeRandomMove(sboard, squares, CLIENTSYMBOL)
	n, err = conn.Write([]byte(sboard))
	if err != nil {
		return n, fmt.Errorf("playTicTacToe first move error while writing %v", sboard)
	}
	fmt.Printf(" S: %q\n", sboard)

	for {
		bytesFromServer := make([]byte, 11)
		n, err = conn.Read(bytesFromServer)
		if err != nil	{
			return n, fmt.Errorf("playTicTacToe() error reading from server %v", err)
		}

		rboard = string(bytesFromServer)
		if strings.Contains(string(bytesFromServer), "END") {
			break
		}

		if !tictactoe.IsValidBoard(rboard)	{
			return n, fmt.Errorf("playTicTacToe() server sent invalid board (%v) %v", rboard, err)
		}
		fmt.Printf(" R: %q\n", rboard)

		if mvCnt, _ := tictactoe.GetMoveDifference(sboard, rboard); mvCnt != 1	{
			return n, fmt.Errorf("playTicTacToe() server made %d moves", mvCnt)
		}

	}

	return 0, nil
}
