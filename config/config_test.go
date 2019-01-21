package config

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	var (
		cfg *Config
	)

	BeforeEach(func() {
		cfg = New()
	})

	Describe("New", func() {
		Context("when instantiating a new config", func() {
			It("should return new config", func() {
				Expect(cfg).ToNot(BeNil())
			})
		})
	})

	Describe("LoadEnvVars", func() {
		Context("when no ENV vars are properly set", func() {
			It("should return an error about the unset vars", func() {
				cfg.Tokens = []string{}

				err := cfg.LoadEnvVars()

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("missing 'GO_MICROSERVICE_1_TOKENS' env var"))
			})
		})

		Context("when all ENV vars are properly set", func() {
			var envVars map[string]string

			BeforeEach(func() {
				envVars = map[string]string{
					"GO_MICROSERVICE_1_TOKENS":         "1111222233334444",
					"GO_MICROSERVICE_1_FOO_API_HOST":   "host",
					"GO_MICROSERVICE_1_MONGO_DB_NAME":  "foo",
					"GO_MICROSERVICE_1_MONGO_DB_HOSTS": "host1",
				}

				for k, v := range envVars {
					os.Setenv(k, v)
				}
			})

			AfterEach(func() {
				for k := range envVars {
					os.Unsetenv(k)
				}
			})

			It("should return nil", func() {
				err := cfg.LoadEnvVars()

				Expect(err).To(BeNil())
			})
		})
	})

	Describe("validateTokens", func() {
		Context("when the passed in tokens are valid", func() {
			It("should return nil", func() {
				tokens := []string{"foooooooooooooooooobar", "biiiiiiinnnnngoooobaaaaannngooo"}
				ok, _ := tokenLength{s: tokens, name: "foo"}.validate()

				Expect(ok).To(BeTrue())
			})
		})

		Context("when the passed in tokens are invalid", func() {
			It("should return nil", func() {
				tokens := []string{"f", "1111222233334444"}
				ok, e := tokenLength{s: tokens, name: "foo"}.validate()

				Expect(ok).ToNot(BeTrue())
				Expect(e).To(ContainElement("f token must be at least 16 chars long"))
			})
		})
	})
})
