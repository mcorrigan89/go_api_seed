# syntax=docker/dockerfile:1

FROM golang:1.22
WORKDIR /usr/src/app

RUN apt update && apt install -y make
RUN apt install -y migrate

COPY . .
RUN go build -o=./bin/main ./cmd

RUN make migrate-up

EXPOSE 9001

CMD ["./bin/main"]
