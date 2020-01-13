package service

import (
	"fmt"
	"github.com/IvanProdaiko94/ssh-test/cfg"
	"github.com/IvanProdaiko94/ssh-test/game"
	"github.com/IvanProdaiko94/ssh-test/persistence"
	"github.com/IvanProdaiko94/ssh-test/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type App struct {
	config *cfg.Config
	db     persistence.TicTacToe
	policy game.Policy
}

func (app *App) GetAPIV1GamesHandler(params operations.GetAPIV1GamesParams) middleware.Responder {
	dbResponse, err := app.db.GetAllGames()
	if err != nil {
		log.Error(err)
		// [404] if failed to find the games
		if err == persistence.ErrNotFound {
			return operations.NewGetAPIV1GamesNotFound()
		}
		// [500] if failed to get the game
		return operations.NewGetAPIV1GamesInternalServerError()
	}
	log.Debug(dbResponse)
	return operations.NewGetAPIV1GamesOK().WithPayload(dbResponse)
}

func (app *App) GetAPIV1GamesGameIDHandler(params operations.GetAPIV1GamesGameIDParams) middleware.Responder {
	dbResponse, err := app.db.GetGameById(string(params.GameID))
	if err != nil {
		log.Error(err)
		// [404] if failed to find the game
		if err == persistence.ErrNotFound {
			return operations.NewGetAPIV1GamesGameIDNotFound()
		}
		// [500] if failed to get the game
		return operations.NewGetAPIV1GamesGameIDInternalServerError()
	}
	log.Debug(dbResponse)
	return operations.NewGetAPIV1GamesGameIDOK().WithPayload(dbResponse)
}

func (app *App) PostAPIV1GamesHandler(params operations.PostAPIV1GamesParams) middleware.Responder {
	// [400] if no param provided
	if params.Game == nil {
		log.Error(errors.New("bad parameters"))
		return operations.NewPostAPIV1GamesBadRequest()
	}

	// [400] if param has wrong length
	if params.Game.Board != nil && len(*params.Game.Board) != 9 {
		log.Error(errors.New("invalid board length"))
		return operations.NewPutAPIV1GamesGameIDBadRequest()
	}

	params.Game.ID = strfmt.UUID(uuid.NewV4().String())
	board := game.NewBoard(*params.Game, app.policy)

	// [400] if param has none of those 3 states:
	// 1 Cross + 8 Free || 1 Nought + 8 Free || 9 Free
	if !board.IsStartOfTheGame() {
		log.Error(errors.New("invalid input"))
		return operations.NewPostAPIV1GamesBadRequest()
	}

	// [500] if machine failed to make a move
	if err := board.MakeMachineMove(game.Noughts); err != nil {
		log.Error(err)
		return operations.NewPostAPIV1GamesInternalServerError()
	}

	model := board.ToModelGame()

	// [500] if failed to create a game
	if err := app.db.CreateGame(model); err != nil {
		log.Error(err)
		return operations.NewPostAPIV1GamesInternalServerError()
	}

	urlBuilder := operations.GetAPIV1GamesGameIDURL{GameID: model.ID}
	url := urlBuilder.String()
	return operations.NewPostAPIV1GamesCreated().
		WithLocation(url).
		WithPayload(&operations.PostAPIV1GamesCreatedBody{Location: url})
}

func (app *App) PutAPIV1GamesGameIDHandler(params operations.PutAPIV1GamesGameIDParams) middleware.Responder {
	// [400] if no param provided
	if params.Game == nil {
		log.Error(errors.New("bad parameters"))
		return operations.NewPutAPIV1GamesGameIDBadRequest()
	}

	// [400] if param has wrong length
	if params.Game.Board != nil && len(*params.Game.Board) != 9 {
		log.Error(errors.New("invalid board length"))
		return operations.NewPutAPIV1GamesGameIDBadRequest()
	}

	prevMove, err := app.db.GetGameById(string(params.Game.ID))
	if err != nil {
		log.Error(errors.Wrap(err, "failed to get previous move"))
	}

	// [400] if board is the same as before or the move was kind of override of the previous one
	if ok := game.ValidateBoardWithPrevMove(*prevMove.Board, *params.Game.Board); !ok {
		log.Error(errors.New(
			fmt.Sprintf("new board %s conflicting with existing one %s", *params.Game.Board, *prevMove.Board),
		))
		return operations.NewPutAPIV1GamesGameIDBadRequest()
	}

	board := game.NewBoard(*params.Game, app.policy)
	// ignore error since if no moves available, than game has come to and end
	_ = board.MakeMachineMove(game.Noughts)
	model := board.ToModelGame()
	if err := app.db.UpdateGame(model); err != nil {
		log.Error(err)
		// [404] if failed to find the game
		if err == persistence.ErrNotFound {
			return operations.NewPutAPIV1GamesGameIDNotFound()
		}
		// [500] if failed to update the game
		return operations.NewPutAPIV1GamesGameIDInternalServerError()
	}
	return operations.NewPutAPIV1GamesGameIDOK().WithPayload(model)
}

func (app *App) DeleteAPIV1GamesGameIDHandler(params operations.DeleteAPIV1GamesGameIDParams) middleware.Responder {
	if err := app.db.DeleteGame(string(params.GameID)); err != nil {
		log.Error(err)
		// [404] if failed to find the game
		if err == persistence.ErrNotFound {
			return operations.NewDeleteAPIV1GamesGameIDNotFound()
		}
		// [500] if failed to delete the game
		return operations.NewDeleteAPIV1GamesGameIDInternalServerError()
	}
	return operations.NewDeleteAPIV1GamesGameIDOK()
}

func (app *App) Close() error {
	return app.db.Close()
}

func New(config *cfg.Config, db persistence.TicTacToe, policy game.Policy) *App {
	return &App{
		config: config,
		db:     db,
		policy: policy,
	}
}
