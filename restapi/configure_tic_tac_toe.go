// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/IvanProdaiko94/ssh-test/restapi/operations"
	"github.com/IvanProdaiko94/ssh-test/service"
	"github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//go:generate swagger generate server --target ../../ssh-test --name TicTacToe --spec ../api/swagger.yaml

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

	api.DeleteAPIV1GamesGameIDHandler = operations.DeleteAPIV1GamesGameIDHandlerFunc(service.App.DeleteAPIV1GamesGameIDHandler)
	api.GetAPIV1GamesHandler = operations.GetAPIV1GamesHandlerFunc(service.App.GetAPIV1GamesHandler)
	api.GetAPIV1GamesGameIDHandler = operations.GetAPIV1GamesGameIDHandlerFunc(service.App.GetAPIV1GamesGameIDHandler)
	api.PostAPIV1GamesHandler = operations.PostAPIV1GamesHandlerFunc(service.App.PostAPIV1GamesHandler)
	api.PutAPIV1GamesGameIDHandler = operations.PutAPIV1GamesGameIDHandlerFunc(service.App.PutAPIV1GamesGameIDHandler)

	api.ServerShutdown = func() {
		err := service.App.Close()
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
	return handler
}
