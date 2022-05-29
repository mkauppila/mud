FROM golang:1.18-alpine3.16 AS build

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY internal/server ./internal/server
COPY cmd/server.go ./cmd/server.go
RUN go build -o server cmd/server.go

FROM alpine:3.16.0 AS final
COPY --from=build /usr/src/app/server .
ENTRYPOINT [ "./server" ]

