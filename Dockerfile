FROM golang:1.14-alpine as dev

WORKDIR /go/src/snowflake
COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .
RUN go install github.com/dustinpianalto/snowflake/...

CMD [ "go", "run", "cmd/snowflake/main.go" ]

FROM alpine

WORKDIR /bin

COPY --from=dev /go/bin/snowflake ./snowflake

CMD [ "snowflake" ]
