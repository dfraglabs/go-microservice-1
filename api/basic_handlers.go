package api

import (
	"net/http"

	"github.com/InVisionApp/rye"
	"github.com/InVisionApp/go-health"
)

type HealthcheckStatus struct {
	Details map[string]health.State `json:"details"`
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
}

// @Summary Greets you with a friendly message
// @Description If this endpoint does not work, something is seriously busted
// @Tags basic
// @Produce json
// @Success 200 {object} rye.JSONStatus "The service was able to start enough to be able to service inbound requests"
// @Router / [get]
func (a *API) homeHandler(rw http.ResponseWriter, r *http.Request) {
	rye.WriteJSONStatus(rw, "Oh, hello there!", "Refer to README.md for dfraglabs/go-microservice-1 API usage", http.StatusOK)
}

// @Summary Returns the current version of the service
// @Description Another simple handler, similar to '/' - if this does not work, something is broken
// @Tags basic
// @Produce json
// @Success 200 {object} rye.JSONStatus "'status' contains the string 'version', while 'message' will contain the actual version"
// @Router /version [get]
func (a *API) versionHandler(rw http.ResponseWriter, r *http.Request) {
	rye.WriteJSONStatus(rw, "version", "dfraglabs/go-microservice-1 "+a.Version, http.StatusOK)
}

// @Summary Describes the current health of fleet-api
// @Description Describes the current health of fleet-api
// @Tags basic
// @Produce json
// @Success 200 {object} api.HealthcheckStatus "All is well"
// @Failure 500 {object} api.HealthcheckStatus "The service is unhealthy"
// @Router /health [get]
func dummyHealth() {}

// @Summary View API docs via Swagger-UI
// @Description This endpoint serves the API spec via Swagger-UI (using github.com/swaggo/swag)
// @Tags basic
// @Produce html
// @Success 200 {string} string "Swagger-UI"
// @Router /docs/index.html [get]
func dummyDocs() {}