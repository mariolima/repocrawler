# Dockerfile References: https://docs.docker.com/engine/reference/builder/

FROM golang:latest

LABEL maintainer="Mario Lima <msclima@uporto.pt>"

WORKDIR /app

RUN ls

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/crawler-cli

RUN go build -o crawler .

EXPOSE 8080

CMD ["./crawler"]
