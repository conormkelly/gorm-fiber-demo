# Go Fiber Demo

## Intro

This Go app is a basic crud API written in a similar style to Express using Fiber.
By no means am I an expert in Go, so take the actual coding style with a pinch of salt.

Having said that, it is intended to demonstrate testing and Dockerized testing approaches:

- [Table-driven](#unit--integration-testing) unit / integration tests in multiple packages (main, services, database) that can be executed via `go test ./...` or other commands described later in the README. There is ~80% code coverage.
- A multi-stage Dockerfile that uses [profiles](https://docs.docker.com/compose/profiles) in the `docker-compose.yml` file. This offers more flexibility in terms of env setup in future and for local tooling requirements.
- E2E tests via dockerized Newman (CLI-based Postman test runner).

## Setup

### Dev environment and workflow

- NOTE: A recent Docker installation is a pre-requisite.

  I'm using `Docker version 20.10.12`.

1. Setting up development env

First, we generate a config file for _Air_, which is similar to Nodemon, enabling hot-reloading.
The resultant file, `.air.toml` is generated in the locla working directory, and is .gitignored.

```sh
docker compose run --rm my-go-api air init
```

Now, if we run the following command to start the containers, we can make changes to our Go API code and see the changes reflected immediately in the terminal:

```sh
docker compose up
```

## Testing

### Application / API testing

The below script will start containers for the services outlined in `docker-compose.yml`:

- Golang API
- MySQL DB
- Newman (Postman test runner)

To execute it, make sure you have Docker installed and run:

```sh
./testing/test.sh
```

This will output the results of the Postman test and bring down the containers.

_Note: No actual collection tests are in place yet._

### Unit / integration testing

The core tests are:

- Integration tests that operate across all layers of the app, using an in-memory SQLite database.
- Focused on testing at the feature level - blackbox and less brittle than the unit tests.
- [Table-driven](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests).

This allows new test cases to be easily added without modifying the existing test methods.

It also utilizes the Go testing feature of [subtests](https://go.dev/blog/subtests), which allows for greater flexibility when using table-driven tests,
such as [running a specific set of subtests](#running-a-specific-set-of-subtests).

### Handy test commands

It may be helpful to add these commands as functions in `~/.bash_profile`.

#### To test all packages

```sh
go test ./...
```

#### To view overall coverage percentage

```sh
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
go tool cover -func profile.cov

```

#### Generating code coverage HTML

```sh
 go test -covermode=set -coverpkg=./... -coverprofile coverage.out -v ./...
 go tool cover -html coverage.out -o coverage.html
```

#### Running a test suite

```sh
go test -run=TestGetUser
```

#### Running a specific set of subtests

```sh
# The contents in quotes are matched against test case description e.g.
# in this case, the "Get all users when table is empty" subtest
go test -run=TestGetAllUsers/"table is empty"
```

## Todo

- Look into validation lib recommended by Fiber team in their docs
- Swagger integration

## Resources

- Great primer on using Docker + Newman together

  <https://www.wwt.com/article/postman-api-tests-collection-run-with-docker-compose>

- Multi-stage Docker compose-build

  <https://firehydrant.com/blog/develop-a-go-app-with-docker-compose>
