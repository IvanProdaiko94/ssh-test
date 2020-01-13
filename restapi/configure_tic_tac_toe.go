// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/IvanProdaiko94/ssh-test/cfg"
	"github.com/IvanProdaiko94/ssh-test/game"
	"github.com/IvanProdaiko94/ssh-test/persistence/postgres"
	"github.com/IvanProdaiko94/ssh-test/restapi/operations"
	"github.com/IvanProdaiko94/ssh-test/service"
	"github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

//go:generate swagger generate server --target ../../ssh-test --name TicTacToe --spec ../api/swagger.yaml

var application *service.App

func init() {
	// FIXME: I failed to find batter place to init all of the stuff
	config := cfg.ReadEnv()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)

	db, err := postgres.InitDBConnection(config.PostgresConfig)
	if err != nil {
		panic(err)
	}
	var policy game.Policy
	if config.PolicyFilePath != "" {
		var err error
		policy, err = game.NewDefaultPolicy(config.PolicyFilePath)
		if err != nil {
			panic(err)
		}
	}
	application = service.New(config, postgres.NewTicTacToe(db), policy)
}

func configureFlags(api *operations.TicTacToeAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.TicTacToeAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.DeleteAPIV1GamesGameIDHandler = operations.DeleteAPIV1GamesGameIDHandlerFunc(application.DeleteAPIV1GamesGameIDHandler)
	api.GetAPIV1GamesHandler = operations.GetAPIV1GamesHandlerFunc(application.GetAPIV1GamesHandler)
	api.GetAPIV1GamesGameIDHandler = operations.GetAPIV1GamesGameIDHandlerFunc(application.GetAPIV1GamesGameIDHandler)
	api.PostAPIV1GamesHandler = operations.PostAPIV1GamesHandlerFunc(application.PostAPIV1GamesHandler)
	api.PutAPIV1GamesGameIDHandler = operations.PutAPIV1GamesGameIDHandlerFunc(application.PutAPIV1GamesGameIDHandler)

	api.ServerShutdown = func() {
		err := application.Close()
		log.Errorf("close sqldb connection error: %s", err.Error())
		// postgres shutdown
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	//handlePanic := recover.New(&recover.Options{Log: log.Println})
	//logViaLogrus := interpose.NegroniLogrus()

	// rate limiter

	//limiter := tollbooth.NewLimiter(10, nil)
	//limiter.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
	// tollbooth.LimitHandler(limiter, logViaLogrus(handler))

	//return handlePanic(logViaLogrus(handler))
	return handler
}
