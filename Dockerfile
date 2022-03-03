FROM golang:1.17-alpine

WORKDIR /go/src
COPY . .

RUN go get -d -v ./...
RUN go install ./main.go

CMD ["main"]