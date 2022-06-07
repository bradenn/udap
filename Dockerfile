FROM golang:1.18rc1-stretch

WORKDIR /udap
COPY . .

RUN go get -d -v ./...
RUN go install ./cmd/udap/main.go

CMD ["main"]