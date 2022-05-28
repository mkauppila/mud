FROM golang:1.18-alpine3.16 AS build

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY mud/ ./mud/
COPY main.go ./
RUN go build -o server

FROM alpine:3.16.0 AS final
COPY --from=build /usr/src/app/server .
ENTRYPOINT [ "./server" ]

