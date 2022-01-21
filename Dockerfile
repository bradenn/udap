FROM golang:1.18beta1-alpine3.15

WORKDIR /go/src
COPY . .

RUN go get -d -v ./...
RUN go install -v ./cmd/run/main.go

CMD ["main"]