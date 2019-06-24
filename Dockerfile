FROM golang:1.12

WORKDIR /go/src/app
COPY . .

RUN go install -v ./...

EXPOSE 8080

ENTRYPOINT ["web"]