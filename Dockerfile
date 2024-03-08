# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /

# Get dependencies
COPY go.mod go.sum ./
RUN go mod download

# Neorgify src
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build

# Neorgify Config
COPY server ./
COPY login ./
COPY dockerfolder ./folder

RUN mkdir ./Folder

CMD ["/Neorgify"]
