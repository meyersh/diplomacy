package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func duplicateInRange(slice string) bool {
	var seen [10]bool
	for _, c := range slice {
		if c == '.' {
			continue
		}
		var n = c - '0' // convert '5' to 5
		if seen[n] {
			return true // duplicate detected so stop checking.
		}

		seen[n] = true
	}

	return false
}

func validateRow(board string, row int) bool {
	// For a given board string, validate row (0-8)
	// 0 = 0->8
	// 1 = 9->17
	// 2 = 18->26
	// ...
	// n = n->n+8

	log.Printf("row: %d, slice: %s\n", row, board[row*9:row*9+9])

	return !duplicateInRange(board[row*9 : row*9+9])
}

func validateCol(board string, col int) bool {
	// For a given board string, validate column (0-8)
	// 0 = 0, 9, 18, ...
	// 1 = 0+1, 9+1, 18+1, ...
	// 2 = 0+2, 9+2, 18+2, ...
	// ...
	// n = n->n+8

	var colval string
	for i := 0; i < 9; i++ {
		colval += string(rune(board[i*9+col]))
	}

	log.Printf("col: %d, slice: %s\n", col, colval)

	return !duplicateInRange(colval)
}

func validateSquare(board string, square int) bool {
	// For a given board string, validate square (0-8)
	/*
	   0 | 1 | 2
	   3 | 4 | 5
	   6 | 7 | 8

	   row = (square / 3) * 3  // 0, 0, 0, 3, 3, 3, 6, 6, 6
	   col = square % 3 * 3    // 0, 3, 6, 0, 3, 6, 0, 3, 6
	*/
	// 0 = 0->2, 9->11, 18->20
	// 1 = 3->5, 12->14, 21->23
	// 3 = 6->8, 15->17, 24->26
	// 4 = 27->29, 36->38, 45->47
	// 5 = 30->32, 39->41, 48->50
	// ...
	// THEREFORE,
	// n = row*9+col -> row*9+col+2, (row+2)*9+col -> (row+2)*9+col+2, (row+3)*9+col -> (row+3)*9+col+2

	var row = (square / 3) * 3 // 0, 0, 0, 3, 3, 3, 6, 6, 6
	var col = (square % 3) * 3 // 0, 3, 6, 0, 3, 6, 0, 3, 6

	// Build a "cell" string by concatonating the 3x3 square components together.
	log.Printf("square %d -> row=%d, col=%d\n", square, row, col)
	log.Printf("       %d->%d, %d->%d, %d->%d\n", row*9+col, row*9+(col+2),
		(row+1)*9+col, (row+1)*9+(col+2), +(row+2)*9+col, (row+2)*9+(col+2))

	var cell string = board[row*9+col:row*9+(col+3)] +
		board[(row+1)*9+col:(row+1)*9+(col+3)] +
		board[(row+2)*9+col:(row+2)*9+(col+3)]

	log.Printf("       %s\n", cell)

	return !duplicateInRange(cell)
}

func sudokuValidator(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		// fmt.Fprintf(w, "Well, you need to POST a board.")
		// fmt.Fprintf(w, "Try `curl -X POST -d 'board=...1.2.3.4....' http://localhost:8090/sudokuValidator")
		return
	}

	if err := req.ParseForm(); err != nil {
		// fmt.Fprintf(w, "ParseForm() err: %v", err)
		http.Error(w, "ParseForm() error", http.StatusMethodNotAllowed)
		return
	}

	board := req.FormValue("board")
	var boardOk bool = true
	var reason string = ""

	if len(board) != 81 {
		boardOk = false
		reason = "Bad board size."
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// fmt.Fprintf(w, "Submitted board='%s'\n", board)

	for i := 0; i < 9; i++ {
		if !validateCol(board, i) {
			boardOk = false
			reason = "Bad column"
			break
		}
		if !validateRow(board, i) {
			boardOk = false
			reason = "Bad row"
			break
		}
		if !validateSquare(board, i) {
			boardOk = false
			reason = "Bad square"
			break
		}
	}

	// fmt.Fprintf(w, "BoardOK? %t, Reason='%s'\n", boardOk, reason)
	if !boardOk {
		http.Error(w, reason, http.StatusMethodNotAllowed)
		return
	}
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/sudokuValidator", sudokuValidator)

	http.ListenAndServe(":8090", nil)
}
