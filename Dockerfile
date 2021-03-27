FROM golang:1.12

WORKDIR /go/src/github.com/olegfomenko/game-service

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/game-service github.com/olegfomenko/game-service


###

FROM alpine:3.9

COPY --from=0 /usr/local/bin/game-service /usr/local/bin/game-service
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["game-service"]
