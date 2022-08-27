# Go Fiber Demo

# Setup

## Hot reloading (optional)

It is optional but recommended to install _Air_, which allows for hot-reloading on file changes for improved developer experience.

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
