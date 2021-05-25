FROM golang:1.16.4-alpine

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["./hr"]