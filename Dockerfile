# "Inspired" by:
# https://firehydrant.com/blog/develop-a-go-app-with-docker-compose/

FROM golang:1.18 as base

FROM base as dev

# Install Air (Go equivalent to nodemon) for hot-reloading
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /opt/app/api
CMD ["air"]
