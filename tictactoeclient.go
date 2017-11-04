package main

import (
	"fmt"
	"github.com/arjunkrishnababu96/tictactoe"
	"log"
	"net"
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
	var rboard string
	var sboard string = tictactoe.GetEmptyBoard()
	var clientWon, serverWon bool

	var n int
	var err error

	// make first move before the infinite loop starts
	sboard, _ = tictactoe.MakeRandomMove(sboard, tictactoe.AllSquares, tictactoe.CLIENTSYMBOL)
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
		if rboard == "END" {
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
			sboard = tictactoe.SERVERWON
			serverWon = true
		} else if win, ptrn := tictactoe.CanWinNext(rboard, tictactoe.CLIENTSYMBOL); win {
			sboard, _ = tictactoe.MakeWinMove(rboard, ptrn, tictactoe.CLIENTSYMBOL)
			clientWon = true
		} else if win, ptrn := tictactoe.CanWinNext(rboard, tictactoe.SERVERSYMBOL); win {
			sboard, _ = tictactoe.BlockWinMove(rboard, ptrn, tictactoe.CLIENTSYMBOL)
		} else {
			sboard, err = tictactoe.MakeRandomMove(rboard, tictactoe.AllSquares, tictactoe.CLIENTSYMBOL)
			if err != nil {
				// no more empty positions
				sboard = tictactoe.TIE
			}
		}

		n, err = conn.Write([]byte(sboard))
		if err != nil {
			return n, fmt.Errorf("playTicTacToe() error while writing %v", sboard)
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
