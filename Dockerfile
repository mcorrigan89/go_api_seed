# syntax=docker/dockerfile:1

FROM golang:1.22
WORKDIR /usr/src/app

RUN apt update && apt install -y make
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN make migrate-up

COPY . .
RUN go build -o=./bin/main ./cmd

EXPOSE 9001

CMD ["./bin/main"]
