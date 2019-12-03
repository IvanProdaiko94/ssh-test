// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/IvanProdaiko94/ssh-test/restapi/operations"
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
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.DeleteAPIV1GamesGameIDHandler == nil {
		api.DeleteAPIV1GamesGameIDHandler = operations.DeleteAPIV1GamesGameIDHandlerFunc(func(params operations.DeleteAPIV1GamesGameIDParams) middleware.Responder {
			fmt.Print("xxx")
			return middleware.NotImplemented("operation .DeleteAPIV1GamesGameID has not yet been implemented")
		})
	}
	if api.GetAPIV1GamesHandler == nil {
		api.GetAPIV1GamesHandler = operations.GetAPIV1GamesHandlerFunc(func(params operations.GetAPIV1GamesParams) middleware.Responder {
			return middleware.NotImplemented("operation .GetAPIV1Games has not yet been implemented")
		})
	}
	if api.GetAPIV1GamesGameIDHandler == nil {
		api.GetAPIV1GamesGameIDHandler = operations.GetAPIV1GamesGameIDHandlerFunc(func(params operations.GetAPIV1GamesGameIDParams) middleware.Responder {
			return middleware.NotImplemented("operation .GetAPIV1GamesGameID has not yet been implemented")
		})
	}
	if api.PostAPIV1GamesHandler == nil {
		api.PostAPIV1GamesHandler = operations.PostAPIV1GamesHandlerFunc(func(params operations.PostAPIV1GamesParams) middleware.Responder {
			return middleware.NotImplemented("operation .PostAPIV1Games has not yet been implemented")
		})
	}
	if api.PutAPIV1GamesGameIDHandler == nil {
		api.PutAPIV1GamesGameIDHandler = operations.PutAPIV1GamesGameIDHandlerFunc(func(params operations.PutAPIV1GamesGameIDParams) middleware.Responder {
			return middleware.NotImplemented("operation .PutAPIV1GamesGameID has not yet been implemented")
		})
	}

	api.ServerShutdown = func() {}

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
