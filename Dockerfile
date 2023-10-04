FROM golang:1.21.0

WORKDIR /app

ADD . /app

COPY .env /app/cmd/.env

WORKDIR /app/cmd

RUN go mod download

RUN go test -v ./...

RUN go build -o main .

CMD ["/app/cmd/main"]