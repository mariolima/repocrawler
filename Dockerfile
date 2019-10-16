# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# FROM node:latest as build-deps

# WORKDIR /app/web

# COPY . ./

# RUN yarn

# RUN yarn build

FROM golang:latest

LABEL maintainer="Mario Lima <msclima@uporto.pt>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/crawler-cli

RUN go build -o crawler .

EXPOSE 8090

ENTRYPOINT ["./crawler"]
