package logic

import (
	"encoding/json"
	"io/ioutil"
)

type Policy interface {
	FindBestMove(board string) int
}

type defaultPolicy struct {
	Moves map[string]int
}

func (p *defaultPolicy) FindBestMove(board string) int {
	return p.Moves[board]
}

func NewDefaultPolicy(policyFilePath string) (Policy, error) {
	defaultPolicy := &defaultPolicy{}
	movesFile, err := ioutil.ReadFile(policyFilePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(movesFile, &defaultPolicy.Moves)
	if err != nil {
		return nil, err
	}
	return defaultPolicy, nil
}
