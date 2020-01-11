package service

import (
	"errors"
	"github.com/IvanProdaiko94/ssh-test/cfg"
	"github.com/IvanProdaiko94/ssh-test/logic"
	"github.com/IvanProdaiko94/ssh-test/models"
	"github.com/IvanProdaiko94/ssh-test/persistence"
	"github.com/IvanProdaiko94/ssh-test/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type App struct {
	config cfg.Config
	api    *operations.TicTacToeAPI
	db     persistence.TicTacToe
}

func (app *App) GetAPIV1GamesHandler(params operations.GetAPIV1GamesParams) middleware.Responder {
	games, err := app.db.GetAllGames()
	if err != nil {
		log.Error(err)
		if err == persistence.ErrNotFound {
			return operations.NewGetAPIV1GamesNotFound()
		}
		return operations.NewGetAPIV1GamesInternalServerError()
	}
	return operations.NewGetAPIV1GamesOK().WithPayload(games)
}

func (app *App) GetAPIV1GamesGameIDHandler(params operations.GetAPIV1GamesGameIDParams) middleware.Responder {
	if len(params.GameID.String()) == 0 {
		log.Error(errors.New("bad parameters"))
		return operations.NewGetAPIV1GamesGameIDBadRequest()
	}
	game, err := app.db.GetGameById(string(params.GameID))
	if err != nil {
		log.Error(err)
		if err == persistence.ErrNotFound {
			return operations.NewGetAPIV1GamesGameIDNotFound()
		}
		return operations.NewGetAPIV1GamesGameIDInternalServerError()
	}
	return operations.NewGetAPIV1GamesGameIDOK().WithPayload(game)
}

func (app *App) PostAPIV1GamesHandler(params operations.PostAPIV1GamesParams) middleware.Responder {
	game := params.Game
	if game == nil {
		log.Error(errors.New("bad parameters"))
		return operations.NewPostAPIV1GamesBadRequest()
	}

	board := logic.NewBoard(*game, app.config.GameMode)
	if status := board.CurrentStatus(); status != models.GameStatusRUNNING {
		log.Error(errors.New("invalid input"))
		return operations.NewPostAPIV1GamesBadRequest()
	}
	// FIXME: not always O
	if err := board.MakeMachineMove(logic.Noughts); err != nil {
		log.Error(err)
		return operations.NewPostAPIV1GamesInternalServerError()
	}

	model := board.ToModel()
	model.ID = strfmt.UUID(uuid.NewV4().String())

	err := app.db.CreateGame(model)
	if err != nil {
		log.Error(err)
		return operations.NewPostAPIV1GamesInternalServerError()
	}

	urlBuilder := operations.GetAPIV1GamesGameIDURL{GameID: model.ID}
	url := urlBuilder.String()
	return operations.NewPostAPIV1GamesCreated().WithLocation(url)
}

func (app *App) PutAPIV1GamesGameIDHandler(params operations.PutAPIV1GamesGameIDParams) middleware.Responder {
	game := params.Game
	if game == nil {
		log.Error(errors.New("bad parameters"))
		return operations.NewPutAPIV1GamesGameIDBadRequest()
	}

	board := logic.NewBoard(*game, app.config.GameMode)
	// FIXME: not always O
	// ignore error since if no moves available, than game has come to and end
	_ = board.MakeMachineMove(logic.Noughts)
	model := board.ToModel()
	err := app.db.UpdateGame(model)
	if err != nil {
		log.Error(err)
		return operations.NewPutAPIV1GamesGameIDInternalServerError()
	}
	return operations.NewPutAPIV1GamesGameIDBadRequest()
}

func (app *App) DeleteAPIV1GamesGameIDHandler(params operations.DeleteAPIV1GamesGameIDParams) middleware.Responder {
	if len(params.GameID.String()) == 0 {
		return operations.NewDeleteAPIV1GamesGameIDBadRequest()
	}
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

func New(config cfg.Config, api *operations.TicTacToeAPI, db persistence.TicTacToe) *App {
	return &App{
		config: config,
		api:    api,
		db:     db,
	}
}
