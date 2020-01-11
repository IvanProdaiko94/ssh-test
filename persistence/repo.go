package persistence

import (
	"github.com/IvanProdaiko94/ssh-test/models"
	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type SQLDBConfig struct {
	Host         string
	Port         int
	User         string
	Pass         string
	DBName       string
	MaxOpenConns int
}

type TicTacToe interface {
	GetAllGames() ([]*models.Game, error)
	GetGameById(id string) (*models.Game, error)
	CreateGame(game *models.Game) error
	UpdateGame(game *models.Game) error
	DeleteGame(id string) error

	Close() error
	Check() error
}
