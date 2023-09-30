FROM golang:1.21.0
WORKDIR /app
ADD . /app
WORKDIR /app/cmd
RUN go mod download
RUN go build -o main .
CMD ["/app/cmd/main"]