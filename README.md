# Go Fiber Demo

# Setup

## Dev Tools

### Hot-reloading via Air

It is optional but highly recommended to install _Air_, which allows for hot-reloading on file changes for improved developer experience.

There are 2 ways to install it.

1. Manually

   ```sh
   # binary will be $(go env GOPATH)/bin/air
   curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

   # or install it into ./bin/
   curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
   ```

2. As a Go module (requires Go 1.16+)

   ```sh
   go install github.com/cosmtrek/air@latest
   ```

Confirm installation by running:

```sh
air -v
```

## Running the app

```sh
# If using air:
air

# Otherwise:
go run .
```

## Testing

To test all packages:

```sh
go test ./...
```

To view overall coverage percentage:

```sh
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
go tool cover -func profile.cov

```

Generating code coverage HTML:

```sh
 go test -covermode=set -coverpkg=./... -coverprofile coverage.out -v ./...
 go tool cover -html coverage.out -o coverage.html
```

It may be helpful to add these commands as functions in `~/.bash_profile`.

## Todo

- Look into validation frameworks (similar to AJV etc)
- Swagger integration
