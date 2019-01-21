package api

import (
	"net/http"
	"os"

	hh "github.com/InVisionApp/go-health/handlers"
	"github.com/InVisionApp/rye"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/swaggo/http-swagger"

	"github.com/dfraglabs/go-microservice-1/config"
	"github.com/dfraglabs/go-microservice-1/deps"
	_ "github.com/dfraglabs/go-microservice-1/docs"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("pkg", "api")
}

type API struct {
	Config  *config.Config
	Version string
	Deps    *deps.Dependencies
}

type APIResponseJSON struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Values  map[string]string `json:"values,omitempty"`
	Errors  string            `json:"errors,omitempty"`
}

func New(cfg *config.Config, d *deps.Dependencies, version string) *API {
	return &API{
		Config:  cfg,
		Version: version,
		Deps:    d,
	}
}

func (a *API) Run() error {
	llog := log.WithField("method", "Run")
	llog.Infof("Starting API server...")

	routes := mux.NewRouter().StrictSlash(true)

	/**************
	 * Basic handlers
	 **************/

	routes.Handle(
		"/", http.HandlerFunc(a.homeHandler),
	).Methods("GET")

	routes.Handle(
		"/version", http.HandlerFunc(a.versionHandler),
	).Methods("GET")

	healthHandler := hh.NewJSONHandlerFunc(a.Deps.Health, map[string]interface{}{
		"version": a.Version,
	})

	routes.Handle(newrelic.WrapHandle(a.Deps.NRApp,
		"/healthcheck", healthHandler,
	)).Methods("GET")

	// Expose API spec via /docs/index.html (requires initial `make docs` run)
	routes.PathPrefix("/docs").HandlerFunc(
		httpSwagger.WrapHandler,
	).Methods("GET")

	/**************
	 *  v1 endpoints
	 **************/

	llog.Infof("API server running on %v", a.Config.ListenAddress)

	return http.ListenAndServe(a.Config.ListenAddress, routes)
}

func (a *API) setupHandler(path string, ryeStack []rye.Handler) (string, http.Handler) {
	p, h := newrelic.WrapHandle(a.Deps.NRApp, path, a.Deps.MWHandler.Handle(ryeStack))
	return p, handlers.LoggingHandler(os.Stdout, h)
}
