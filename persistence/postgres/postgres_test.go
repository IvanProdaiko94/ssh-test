package postgres

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanProdaiko94/ssh-test/models"
	"github.com/IvanProdaiko94/ssh-test/persistence"
	"github.com/go-openapi/strfmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
)

var games = []*models.Game{
	{
		Board:  new(string),
		ID:     strfmt.UUID("1"),
		Status: models.GameStatusRUNNING,
	},
	{
		Board:  new(string),
		ID:     strfmt.UUID("2"),
		Status: models.GameStatusDRAW,
	},
	{
		Board:  new(string),
		ID:     strfmt.UUID("3"),
		Status: models.GameStatusXWON,
	},
}

type PostgresSuite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository persistence.TicTacToe
}

func (s *PostgresSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)

	s.DB.LogMode(true)

	s.repository = NewTicTacToe(s.DB)
	require.NoError(s.T(), err)
}

func (s *PostgresSuite) Test_repository_GetAllGames() {
	rows := sqlmock.NewRows([]string{"board", "id", "status"})
	for _, game := range games {
		rows.AddRow(game.Board, game.ID, game.Status)
	}

	s.mock.ExpectQuery(
		regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM "%s"`, GamesTable)),
	).WillReturnRows(rows)

	result, err := s.repository.GetAllGames()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), games, result)
}

func (s *PostgresSuite) Test_repository_GetGameById() {
	rows := sqlmock.NewRows([]string{"board", "id", "status"})
	for _, game := range games {
		rows.AddRow(game.Board, game.ID, game.Status)
	}
	s.mock.ExpectQuery(
		regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM "%s" WHERE (id = $1)`, GamesTable)),
	).
		WithArgs(games[0].ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"board", "id", "status"}).AddRow(games[0].Board, games[0].ID, games[0].Status),
		)

	result, err := s.repository.GetGameById(string(games[0].ID))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), games[0], result)
}

func (s *PostgresSuite) Test_repository_CreateGame() {
	rowsAffected := driver.RowsAffected(int64(1))
	s.mock.
		ExpectExec(fmt.Sprintf(
			"CREATE %s SET board = $1, id = $2, status = $3",
			GamesTable,
		)).
		WithArgs(games[0].Board, games[0].ID, games[0].Status).
		WillReturnResult(rowsAffected)

	err := s.repository.CreateGame(games[0])
	assert.NoError(s.T(), err)
}

func (s *PostgresSuite) Test_repository_UpdateGame() {
	rowsAffected := driver.RowsAffected(int64(1))
	s.mock.
		ExpectExec(fmt.Sprintf(
			"UPDATE %s SET board = $1, id = $2, status = $3 WHERE (id = $4)",
			GamesTable,
		)).
		WithArgs(games[0].Board, games[0].ID, games[0].Status, games[0].ID).
		WillReturnResult(rowsAffected)

	err := s.repository.UpdateGame(games[0])
	assert.NoError(s.T(), err)
}

func (s *PostgresSuite) Test_repository_DeleteGame() {
	rowsAffected := driver.RowsAffected(int64(1))
	s.mock.
		ExpectExec(fmt.Sprintf(
			"DELETE * FROM %s WHERE (id = $4)",
			GamesTable,
		)).
		WithArgs(games[0].ID).
		WillReturnResult(rowsAffected)

	err := s.repository.DeleteGame(string(games[0].ID))
	assert.NoError(s.T(), err)
}

func (s *PostgresSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}
