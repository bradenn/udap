FROM golang:1.18-stretch

WORKDIR /udap
COPY . .

RUN go get -d -v ./...
RUN go get -d -v ./modules/**

RUN which go

RUN go install ./cmd/udap/main.go

CMD ["main"]