package game

import (
	"errors"
	"github.com/IvanProdaiko94/ssh-test/models"
	"github.com/go-openapi/strfmt"
	"math/rand"
	"strings"
)

var (
	EOG         = errors.New("game has come to an end")
	Occupied    = errors.New("the cell is occupied already")
	NoFreeMoves = errors.New("no free moves left")
)

type Cell = uint8

const (
	f Cell = iota
	x
	o
)

const (
	Crosses = x
	Noughts = o
)

type Board struct {
	ID     strfmt.UUID
	B      []Cell
	Policy Policy
}

func (gb *Board) String() string {
	board := make([]byte, len(gb.B))
	for i, v := range gb.B {
		switch v {
		case x:
			board[i] = 'x'
		case o:
			board[i] = 'o'
		case f:
			board[i] = '-'
		}
	}
	return string(board)
}

func (gb *Board) hasFreeMoves() bool {
	for _, cell := range gb.B {
		if cell == f {
			return true
		}
	}
	return false
}

func (gb *Board) checkWinner(wc Cell) bool {
	b := gb.B
	// check horizontally
	return (b[0] == wc && b[1] == wc && b[2] == wc) ||
		(b[3] == wc && b[4] == wc && b[5] == wc) ||
		(b[6] == wc && b[7] == wc && b[8] == wc) ||
		// check vertically
		(b[0] == wc && b[3] == wc && b[6] == wc) ||
		(b[1] == wc && b[4] == wc && b[7] == wc) ||
		(b[2] == wc && b[5] == wc && b[8] == wc) ||
		// check diagonals
		(b[0] == wc && b[4] == wc && b[8] == wc) ||
		(b[2] == wc && b[4] == wc && b[7] == wc)
}

func (gb *Board) IsStartOfTheGame() bool {
	nonFreeEntriesCount := 0
	for _, cell := range gb.B {
		if cell != f {
			nonFreeEntriesCount += 1
		}
	}
	return nonFreeEntriesCount <= 1
}

func (gb *Board) IsEmpty() bool {
	for _, cell := range gb.B {
		if cell != f {
			return false
		}
	}
	return true
}

func (gb *Board) MakeMachineMove(player Cell) error {
	if !gb.hasFreeMoves() {
		return NoFreeMoves
	}
	if gb.Policy == nil {
		// stupid player
		min := 0
		max := len(gb.B)
		for {
			i := rand.Intn(max-min+1) + min
			if gb.B[i] == f {
				gb.B[i] = player
				return nil
			}
		}
	} else {
		// player with policy
		i := gb.Policy.FindBestMove(gb.String())
		gb.B[i] = player
	}
	return nil
}

func (gb *Board) GetCurrentStatus() string {
	if gb.checkWinner(x) {
		return models.GameStatusXWON
	}
	if gb.checkWinner(o) {
		return models.GameStatusOWON
	}
	if gb.hasFreeMoves() {
		return models.GameStatusRUNNING
	}
	return models.GameStatusDRAW
}

func (gb *Board) ToModelGame() *models.Game {
	strBoard := strings.ToUpper(gb.String())
	return &models.Game{
		Board:  &strBoard,
		ID:     gb.ID,
		Status: gb.GetCurrentStatus(),
	}
}

func NewBoard(game models.Game, policy Policy) *Board {
	var b []Cell
	if game.Board == nil {
		b = []Cell{
			f, f, f,
			f, f, f,
			f, f, f,
		}
	} else {
		b = make([]Cell, len(*game.Board))
		for i, char := range strings.ToLower(*game.Board) {
			switch char {
			case 'x':
				b[i] = x
			case 'o':
				b[i] = o
			case '-':
				b[i] = f
			}
		}
	}

	return &Board{
		ID:     game.ID,
		B:      b,
		Policy: policy,
	}
}
