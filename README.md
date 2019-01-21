# go-microservice-1

A barebones Golang microservice template.

## Libs

* Router: [gorilla](https://github.com/gorilla/)
    * Middleware: [rye](https://github.com/InVisionApp/rye)
* CLI: [kingpin](https://github.com/alecthomas/kingpin)
* Logging: [logrus](https://github.com/sirupsen/logrus)
    * Logging shim via [go-logger](https://github.com/InVisionApp/go-logger/)
* Health: [go-health](https://github.com/InVisionApp/go-health)
* DB (Mongo): [mgo](gopkg.in/mgo.v2)
* Testing: [ginkgo](github.com/onsi/ginkgo)
* Vendor: [govendor](https://github.com/kardianos/govendor) 

## Batteries Included

* Sample DAL / DAO
* Dependency instantiation pattern via `deps/deps.go`
* Env var fetching & validation
* Docker-ready
    * Separate build from runtime (`Dockerfile` vs `Dockerfile.build`)
    * Run via `docker-compose` (if you like)
* Fleshed out `Makefile`
* CI-ready (for codeship)

## Usage

1. `git clone` this repo
2. Create a new repo for your service
3. `cp -prf` everything except `.git` and `codeship-dockercfg.encrypted` to your new service repo dir
4. Find and replace all occurrences of `go-microservice-1` and `GO_MICROSERVICE_1` with your service name
5. Find and replace all occurrences of `dfraglabs` with your org name
6. Try to run `make test` and `make run`
7. Good luck!

## Documentation
This template uses [swag](https://github.com/swaggo/swag) to generate API documentation.

Run `make docs` to generate swagger api spec. The generated spec is available as
[raw files](./docs/swagger/) and viewable via Swagger-UI on the `/docs/index.html`
endpoint on the service.