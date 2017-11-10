package main

import (
	"fmt"
	"github.com/arjunkrishnababu96/tictactoe"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	addr := os.Args[1]
	fmt.Println("Connecting to", addr)
	conn, err := net.Dial("tcp", addr)
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
	var rboard string
	var sboard string = tictactoe.GetEmptyBoard()
	var clientWon, serverWon bool

	var n int
	var err error

	// make first move before the infinite loop starts
	// sboard, _ = tictactoe.MakeRandomMove(sboard, tictactoe.AllSquares, tictactoe.CLIENTSYMBOL)
	sboard, _ = tictactoe.MakeMove(sboard, 5, tictactoe.CLIENTSYMBOL)
	n, err = conn.Write([]byte(sboard))
	if err != nil {
		return n, fmt.Errorf("playTicTacToe first move error while writing %v", sboard)
	}
	fmt.Printf(" S: %q\n", sboard)

InfiniteLoop:
	for {
		bytesFromServer := make([]byte, 11)
		n, err = conn.Read(bytesFromServer)
		if err != nil {
			return n, fmt.Errorf("playTicTacToe() error reading from server %v", err)
		}

		rboard = string(bytesFromServer)
		if strings.Contains(rboard, tictactoe.TIE) {
			fmt.Println("tie")
			break
		}

		if !tictactoe.IsValidBoard(rboard) {
			return n, fmt.Errorf("playTicTacToe() server sent invalid board (%v) %v", rboard, err)
		}
		fmt.Printf(" R: %q\n", rboard)

		if mvCnt, _ := tictactoe.GetMoveDifference(sboard, rboard); mvCnt != 1 {
			return n, fmt.Errorf("playTicTacToe() server made %d moves", mvCnt)
		}

		if tictactoe.HasWon(rboard, tictactoe.SERVERSYMBOL) {
			// fmt.Println("Server won")
			sboard = tictactoe.SERVERWON
			serverWon = true
		} else if win, ptrn := tictactoe.CanWinNext(rboard, tictactoe.CLIENTSYMBOL); win {
			// fmt.Println("Client can win next")
			sboard, _ = tictactoe.MakeWinMove(rboard, ptrn, tictactoe.CLIENTSYMBOL)
			clientWon = true
		} else if win, ptrn := tictactoe.CanWinNext(rboard, tictactoe.SERVERSYMBOL); win {
			// fmt.Println("Server can win next")
			sboard, _ = tictactoe.BlockWinMove(rboard, ptrn, tictactoe.CLIENTSYMBOL)
		} else if tictactoe.IsFree(rboard, 5) {
			// can play center
			sboard, _ = tictactoe.MakeMove(rboard, 5, tictactoe.CLIENTSYMBOL)
			// fmt.Println("playing center! %v %v", rboard, sboard)

			// DOWN: Play opposite corner
		} else if rboard[0] == tictactoe.SERVERSYMBOL && tictactoe.IsFree(rboard, 10) {
			sboard, _ = tictactoe.MakeMove(rboard, 10, tictactoe.CLIENTSYMBOL)
		} else if rboard[2] == tictactoe.SERVERSYMBOL && tictactoe.IsFree(rboard, 8) {
			sboard, _ = tictactoe.MakeMove(rboard, 8, tictactoe.CLIENTSYMBOL)
		} else if rboard[8] == tictactoe.SERVERSYMBOL && tictactoe.IsFree(rboard, 2) {
			sboard, _ = tictactoe.MakeMove(rboard, 2, tictactoe.CLIENTSYMBOL)
		} else if rboard[10] == tictactoe.SERVERSYMBOL && tictactoe.IsFree(rboard, 0) {
			sboard, _ = tictactoe.MakeMove(rboard, 0, tictactoe.CLIENTSYMBOL)

			// DOWN: Play empty corner
		} else if tictactoe.IsFree(rboard, 0) {
			sboard, _ = tictactoe.MakeMove(rboard, 0, tictactoe.CLIENTSYMBOL)
		} else if tictactoe.IsFree(rboard, 2) {
			sboard, _ = tictactoe.MakeMove(rboard, 2, tictactoe.CLIENTSYMBOL)
		} else if tictactoe.IsFree(rboard, 8) {
			sboard, _ = tictactoe.MakeMove(rboard, 8, tictactoe.CLIENTSYMBOL)
		} else if tictactoe.IsFree(rboard, 10) {
			sboard, _ = tictactoe.MakeMove(rboard, 10, tictactoe.CLIENTSYMBOL)

		} else {
			sboard, err = tictactoe.MakeRandomMove(rboard, tictactoe.AllSquares, tictactoe.CLIENTSYMBOL)
			// fmt.Println("playing random! %v %v", rboard, sboard)
			if err != nil {
				// no more empty positions
				sboard = tictactoe.TIE
			}
			// fmt.Printf(" Random %v %v\n", rboard, sboard)
		}

		n, err = conn.Write([]byte(sboard))
		if err != nil {
			return n, fmt.Errorf("playTicTacToe() error while writing %v", sboard)
		}

		if !tictactoe.IsAnyFree(sboard)	{
			// checks for tie after client made a move
			sboard = tictactoe.TIE
			conn.Write([]byte(sboard))
		}

		switch {
		case sboard == tictactoe.TIE:
			fmt.Println(tictactoe.TIE)
			break InfiniteLoop
		case serverWon:
			fmt.Println(tictactoe.SERVERWON)
			break InfiniteLoop
		case clientWon:
			fmt.Println(tictactoe.CLIENTWON)
			break InfiniteLoop
		default:
			fmt.Printf(" S: %q\n", sboard)
		}
	}

	return 0, nil
}
