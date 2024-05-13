FROM golang:1.21

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o consulate

EXPOSE 8080

HEALTHCHECK --interval=1m --timeout=5s \
  CMD curl -f http://127.0.0.1:8080/ping || exit 1

CMD ["./consulate"]
