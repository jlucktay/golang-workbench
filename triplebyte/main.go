package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type space struct {
	t string
}

func (s *space) String() string {
	return fmt.Sprint(s.t)
}

type board struct {
	board [3][3]space
}

func (b *board) init() {
	b.board = [3][3]space{
		{
			space{t: "-"},
			space{t: "-"},
			space{t: "-"},
		},
		{
			space{t: "-"},
			space{t: "-"},
			space{t: "-"},
		},
		{
			space{t: "-"},
			space{t: "-"},
			space{t: "-"},
		},
	}
}

func (b *board) add(x, y int, token string) {
	if b.board[x][y].t == "-" {
		b.board[x][y].t = token
	}
}

func (b *board) String() string {
	output := ""

	for _, ySpaces := range b.board {
		for xIndex, xSpace := range ySpaces {
			output += fmt.Sprint(xSpace.String())

			if xIndex+1 < len(ySpaces) {
				output += fmt.Sprint("|")
			} else {
				output += fmt.Sprintln()
			}
		}
	}

	return output
}

func (b *board) isFull() bool {
	for y, _ := range b.board {
		for x, _ := range b.board[y] {
			if b.board[y][x].t == "-" {
				return false
			}
		}
	}

	return true
}

func (b *board) aiMove() error {
	for y, _ := range b.board {
		for x, _ := range b.board[y] {
			if b.board[y][x].t == "-" {
				b.board[y][x].t = "X"
				return nil
			}
		}
	}

	return errors.New("no available spaces")
}

func main() {
	// init
	b := board{}
	b.init()
	scanner := bufio.NewScanner(os.Stdin)

	// b.add(1, 1, "X")
	fmt.Printf("%s\n", b.String())
	fmt.Println(b.isFull())
	// b.aiMove()
	// fmt.Printf("%s\n", b.String())

	// play loop
	for scanner.Scan() {
		fmt.Print("Enter X and Y co-ords to play a 'O' e.g. '1 2': ")
		fmt.Println(scanner.Text())
		splitInput := strings.Split(scanner.Text(), " ")
		xIn, _ := strconv.Atoi(splitInput[0]) // TODO: error checking
		yIn, _ := strconv.Atoi(splitInput[1])
		b.add(xIn, yIn, "O")
	}

	if errScan := scanner.Err(); errScan != nil {
		log.Fatal(errScan)
	}
}
