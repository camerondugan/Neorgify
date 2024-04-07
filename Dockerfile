FROM golang:latest AS build-stage
WORKDIR /

# Get dependencies
COPY go.mod go.sum ./
RUN go mod download & mkdir ./Folder

# Neorgify src
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build

FROM gcr.io/distroless/static-debian12:latest AS release-stage
WORKDIR /

COPY --from=build-stage /Neorgify /Neorgify
COPY --from=build-stage /Folder /Folder

# Neorgify Config
COPY dockerfolder ./folder

CMD ["/Neorgify"]
