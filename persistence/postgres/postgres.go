package postgres

import (
	"fmt"
	"github.com/IvanProdaiko94/ssh-test/models"
	"github.com/IvanProdaiko94/ssh-test/persistence"
	"github.com/go-openapi/strfmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" //nolint
	"github.com/pkg/errors"
)

const GamesTable = "games"

func checkError(db *gorm.DB) error {
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return persistence.ErrNotFound
		}
		return db.Error
	}
	if db.RowsAffected == 0 {
		return persistence.ErrNotFound
	}
	return nil
}

func InitDBConnection(c persistence.SQLDBConfig) (*gorm.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Pass, c.Host, c.Port, c.DBName)
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db")
	}

	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	return db, nil
}

type repo struct {
	db *gorm.DB
}

func (r *repo) GetAllGames() ([]*models.Game, error) {
	var games []*models.Game
	if err := checkError(r.db.Table(GamesTable).Find(&games)); err != nil {
		return nil, err
	}
	return games, nil
}

func (r *repo) GetGameById(id string) (*models.Game, error) {
	var game models.Game
	if err := checkError(r.db.Table(GamesTable).Find(&game, "id = ?", id)); err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *repo) CreateGame(game *models.Game) error {
	return checkError(r.db.Table(GamesTable).Create(game))
}

func (r *repo) UpdateGame(game *models.Game) error {
	return checkError(r.db.Table(GamesTable).Where("id = ?", game.ID).Update(game))
}

func (r *repo) DeleteGame(id string) error {
	return checkError(r.db.Table(GamesTable).Delete(&models.Game{ID: strfmt.UUID(id)}))
}

func (r *repo) Close() error {
	return r.db.Close()
}

func (r *repo) Check() error {
	return r.db.DB().Ping()
}

func NewTicTacToe(db *gorm.DB) persistence.TicTacToe {
	return &repo{db: db}
}
