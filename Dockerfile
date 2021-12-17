FROM golang:1.17.3

WORKDIR /go/src
COPY . .

RUN go get -d -v ./...
RUN go install -v ./cmd/run/main.go

CMD ["main"]