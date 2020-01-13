package game

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestValidateBoardWithPrevMove(t *testing.T) {
	prev := "xxxooo---"
	curr := "xxxooo---"
	isValid := ValidateBoardWithPrevMove(prev, curr)
	assert.Equal(t, isValid, false, "Invalid")
}

func TestValidateBoardWithPrevMove2(t *testing.T) {
	prev := "xxxooo---"
	curr := "xxx-oo---"
	isValid := ValidateBoardWithPrevMove(prev, curr)
	assert.Equal(t, isValid, false, "Invalid")
}

func TestValidateBoardWithPrevMove3(t *testing.T) {
	prev := "xxxooo---"
	curr := "xxxoooxxx"
	isValid := ValidateBoardWithPrevMove(prev, curr)
	assert.Equal(t, isValid, false, "Invalid")
}

func TestValidateBoardWithPrevMove4(t *testing.T) {
	prev := "X---O----"
	curr := "XX--O----"
	isValid := ValidateBoardWithPrevMove(prev, curr)
	assert.Equal(t, isValid, true, "Valid")
}
