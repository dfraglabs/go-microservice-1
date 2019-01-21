// API documentation generated via `make docs`
//
// @title go-microservice-1
// @version 1.0
// @description <update description>
// @contact.name <update contact name>
//
package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/dfraglabs/go-microservice-1/api"
	"github.com/dfraglabs/go-microservice-1/config"
	"github.com/dfraglabs/go-microservice-1/deps"
)

var (
	version = "No version specified"

	envFile = kingpin.Flag("envfile", "Local Env file to read at startup").Short('e').Default(".env").String()
	debug   = kingpin.Flag("debug", "Enable debug output").Short('d').Bool()
)

func init() {
	// Parse CLI stuff
	kingpin.Version(version)
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.Parse()
}

func main() {
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// JSON formatter for log output if not running in a TTY
	// because Loggly likes JSON but humans like colors
	if !terminal.IsTerminal(int(os.Stderr.Fd())) {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	llog := logrus.WithField("method", "main")

	llog.WithField("filename", *envFile).Debug("Loading env file")
	if err := godotenv.Load(*envFile); err != nil {
		llog.WithFields(logrus.Fields{"filename": *envFile, "err": err.Error()}).Warn("Unable to load dotenv file")
	}

	cfg := config.New()
	if err := cfg.LoadEnvVars(); err != nil {
		llog.WithError(err).Fatal("Could not instantiate configuration")
	}

	llog = llog.WithField("environment", cfg.EnvName)

	llog.Info("Launching go-microservice-1 API")

	d, err := deps.New(cfg)
	if err != nil {
		llog.WithError(err).Fatal("Could not setup dependencies")
	}

	// Start the API server
	a := api.New(cfg, d, version)
	llog.Fatal(a.Run())
}
