package service

import (
	"errors"
	"github.com/IvanProdaiko94/ssh-test/cfg"
	"github.com/IvanProdaiko94/ssh-test/game"
	"github.com/IvanProdaiko94/ssh-test/models"
	"github.com/IvanProdaiko94/ssh-test/persistence"
	"github.com/IvanProdaiko94/ssh-test/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
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
		if err == persistence.ErrNotFound {
			return operations.NewGetAPIV1GamesNotFound()
		}
		return operations.NewGetAPIV1GamesInternalServerError()
	}
	log.Debug(dbResponse)
	return operations.NewGetAPIV1GamesOK().WithPayload(dbResponse)
}

func (app *App) GetAPIV1GamesGameIDHandler(params operations.GetAPIV1GamesGameIDParams) middleware.Responder {
	if len(params.GameID.String()) == 0 {
		log.Error(errors.New("bad parameters"))
		return operations.NewGetAPIV1GamesGameIDBadRequest()
	}
	dbResponse, err := app.db.GetGameById(string(params.GameID))
	if err != nil {
		log.Error(err)
		if err == persistence.ErrNotFound {
			return operations.NewGetAPIV1GamesGameIDNotFound()
		}
		return operations.NewGetAPIV1GamesGameIDInternalServerError()
	}
	log.Debug(dbResponse)
	return operations.NewGetAPIV1GamesGameIDOK().WithPayload(dbResponse)
}

func (app *App) PostAPIV1GamesHandler(params operations.PostAPIV1GamesParams) middleware.Responder {
	if params.Game == nil {
		log.Error(errors.New("bad parameters"))
		return operations.NewPostAPIV1GamesBadRequest()
	}
	params.Game.ID = strfmt.UUID(uuid.NewV4().String())

	board := game.NewBoard(*params.Game, app.policy)
	if status := board.GetCurrentStatus(); status != models.GameStatusRUNNING {
		log.Error(errors.New("invalid input"))
		return operations.NewPostAPIV1GamesBadRequest()
	}
	// FIXME: not always O
	if err := board.MakeMachineMove(game.Noughts); err != nil {
		log.Error(err)
		return operations.NewPostAPIV1GamesInternalServerError()
	}
	log.Debug(board)
	model := board.ToModelGame()
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
	if params.Game == nil {
		log.Error(errors.New("bad parameters"))
		return operations.NewPutAPIV1GamesGameIDBadRequest()
	}

	board := game.NewBoard(*params.Game, app.policy)
	// FIXME: not always O
	// ignore error since if no moves available, than game has come to and end
	_ = board.MakeMachineMove(game.Noughts)

	log.Debug(board)
	model := board.ToModelGame()
	if err := app.db.UpdateGame(model); err != nil {
		log.Error(err)
		return operations.NewPutAPIV1GamesGameIDInternalServerError()
	}
	return operations.NewPutAPIV1GamesGameIDOK().WithPayload(model)
}

func (app *App) DeleteAPIV1GamesGameIDHandler(params operations.DeleteAPIV1GamesGameIDParams) middleware.Responder {
	if len(params.GameID.String()) == 0 {
		return operations.NewDeleteAPIV1GamesGameIDBadRequest()
	}
	log.Debug(params.GameID)
	if err := app.db.DeleteGame(string(params.GameID)); err != nil {
		log.Error(err)
		if err == persistence.ErrNotFound {
			return operations.NewDeleteAPIV1GamesGameIDNotFound()
		}
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
