FROM golang:1.19.1-alpine3.16 AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main

FROM mongo:5.0.12

COPY --from=build /app .

