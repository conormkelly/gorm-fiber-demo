# Go Fiber Demo

## Setup

### Pre-requisites

- Golang (>= v1.18) installed

  ```sh
  # Install project dependencies
  go get ./...
  ```

- Docker installed

## Running the app

```sh
go run .

# or via Docker (will also spin up DB, run postman tests)
docker-compose up
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

- TODO: Multi-stage Docker compose-build

  <https://firehydrant.com/blog/develop-a-go-app-with-docker-compose>
