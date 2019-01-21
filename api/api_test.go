package api

import (
	"net/http"
	"net/http/httptest"

	"github.com/cactus/go-statsd-client/statsd"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dfraglabs/go-microservice-1/config"
	"github.com/dfraglabs/go-microservice-1/deps"
	"github.com/dfraglabs/go-microservice-1/deps/backends"
)

var _ = Describe("API", func() {
	var (
		request  *http.Request
		response *httptest.ResponseRecorder

		cfg              *config.Config
		d                *deps.Dependencies
		api              *API
		fakeStatsDClient statsd.Statter

		testVersion = "1.0.1"
		testTokens  = []string{"abcdefgh12345678"}
	)

	BeforeEach(func() {
		// Instantiate our fakes
		fakeStatsDClient, _ = statsd.NewNoopClient()

		// Create our fake dependencies
		d = &deps.Dependencies{
			Backends: &backends.Backends{},
			StatsD:   fakeStatsDClient,
		}

		cfg = config.New()
		cfg.Tokens = testTokens

		api = New(cfg, d, testVersion)

		response = httptest.NewRecorder()
	})

	Describe("New", func() {
		Context("when instantiating an api", func() {
			It("should have correct attributes", func() {
				Expect(api.Config).ToNot(BeNil())
				Expect(api.Version).To(Equal(testVersion))
			})
		})
	})

	Describe("HomeHandler", func() {
		Context("when the request is successful", func() {
			It("should return info about API usage", func() {
				api.homeHandler(response, request)
				Expect(response.Code).To(Equal(200))
				Expect(response.Body).To(ContainSubstring("README.md"))
			})
		})
	})

	Describe("VersionHandler", func() {
		Context("when the request is successful", func() {
			It("should return the API version", func() {
				api.versionHandler(response, request)
				Expect(response.Code).To(Equal(200))
				Expect(response.Body).To(ContainSubstring(testVersion))
			})
		})
	})
})
