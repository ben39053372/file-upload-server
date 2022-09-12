FROM golang:1.19.1-alpine3.16 AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main

FROM alpine 

WORKDIR /app

COPY --from=build /app/public ./public
COPY --from=build /app/main .

ENTRYPOINT [ "./main" ]

