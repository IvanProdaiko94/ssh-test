package game

import (
	"github.com/IvanProdaiko94/ssh-test/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBoard(t *testing.T) {
	b := "X-OXOXO--"
	model := models.Game{Board: &b, ID: "123", Status: models.GameStatusRUNNING}
	policy, err := NewDefaultPolicy("./policy.json")
	assert.NoError(t, err)
	board := NewBoard(model, policy)
	board.MakeMachineMove(Noughts)
	assert.Equal(t, models.GameStatusOWON, board.GetCurrentStatus())
	assert.NoError(t, err)
}
