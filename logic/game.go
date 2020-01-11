package logic

import (
	"errors"
	"github.com/IvanProdaiko94/ssh-test/models"
	"github.com/go-openapi/strfmt"
	"math/rand"
)

var (
	EOG         = errors.New("logic has come to an end")
	Occupied    = errors.New("the cell is occupied already")
	NoFreeMoves = errors.New("no free moves left")
)

type Cell uint8

const (
	f Cell = iota
	X
	O
)

type Player = Cell

const (
	Noughts = O
	Crosses = X
)

type Mode int8

const (
	stupid Mode = iota
	model

	ModeStupid = "STUPID"
	ModeModel  = "MODEL"
)

func ModeFromString(s string) Mode {
	if s == ModeModel {
		return model
	}
	return stupid
}

type Board struct {
	ID   strfmt.UUID
	B    []Cell
	Mode Mode
	//MachinePlays Player
}

//func machinePlaysWith(b []Cell) Player {
//	notches := 0
//	crosses := 0
//
//	for _, cell := range b {
//		if cell == X {
//			crosses += 1
//		} else if cell == O {
//			notches += 1
//		}
//	}
//
//	machinePlays := Noughts
//	if notches > crosses {
//		machinePlays = Crosses
//	}
//	return machinePlays
//}

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
	if gb.Mode == stupid {
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

	}
	return nil
}

func (gb *Board) MakeHumanMove(player Cell, i int) error {
	if gb.CurrentStatus() != models.GameStatusRUNNING {
		return EOG
	}
	if gb.B[i] != f {
		return Occupied
	}
	// do move
	gb.B[i] = player
	return nil
}

func (gb *Board) CurrentStatus() string {
	if gb.checkWinner(X) {
		return models.GameStatusXWON
	}
	if gb.checkWinner(O) {
		return models.GameStatusOWON
	}
	if gb.hasFreeMoves() {
		return models.GameStatusRUNNING
	}
	return models.GameStatusDRAW
}

func (gb *Board) ToModel() *models.Game {
	board := make([]byte, len(gb.B))
	for i, v := range gb.B {
		board[i] = byte(v)
	}
	strBoard := string(board)
	return &models.Game{
		Board:  &strBoard,
		ID:     gb.ID,
		Status: gb.CurrentStatus(),
	}
}

func NewBoard(game models.Game, mode Mode) *Board {
	var b []Cell
	if game.Board == nil {
		b = []Cell{
			f, f, f,
			f, f, f,
			f, f, f,
		}
	} else {
		b = make([]Cell, len(*game.Board))
		for i, char := range *game.Board {
			b[i] = Cell(char)
		}
	}

	return &Board{
		ID:   game.ID,
		B:    b,
		Mode: mode,
		//MachinePlays: machinePlaysWith(b),
	}
}
