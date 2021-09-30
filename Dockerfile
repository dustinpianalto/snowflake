FROM golang:alpine as dev

WORKDIR /go/src/snowflake
COPY ./go.mod .
COPY ./go.sum .

RUN go mod download
RUN apk update

RUN apk add protobuf
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
COPY . .
RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative snowflake/snowflake.proto
RUN go install github.com/dustinpianalto/snowflake/...

CMD [ "go", "run", "cmd/snowflake/main.go" ]

FROM alpine

WORKDIR /bin

COPY --from=dev /go/bin/snowflake ./snowflake

CMD [ "/bin/snowflake" ]
