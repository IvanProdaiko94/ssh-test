package game

import (
	"strings"
)

func ValidateBoardWithPrevMove(prev string, curr string) bool {
	prev = strings.ToLower(prev)
	curr = strings.ToLower(curr)
	if prev == curr {
		return false
	}

	f := uint8('-')
	x := uint8('x')
	o := uint8('o')

	thereWereChanges := false
	for i := 0; i < len(prev); i++ {
		prevChar := prev[i]
		currChar := curr[i]
		// the only correct move is when
		// on previous step the position was free
		// and now it is occupied
		if prev[i] == f && (currChar == x || currChar == o) {
			// there must be only one valid change at a time
			if thereWereChanges {
				return false
			}
			thereWereChanges = true
			continue
		}
		if prevChar != curr[i] {
			return false
		}
	}
	return true
}
