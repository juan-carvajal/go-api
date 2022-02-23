FROM golang:1.16 as base

FROM base as dev

RUN go get -u github.com/cosmtrek/air

WORKDIR /app
ENTRYPOINT ["air"]