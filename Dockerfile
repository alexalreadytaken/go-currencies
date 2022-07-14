FROM golang:alpine as builder

WORKDIR /curr-build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o curr cmd/main.go

FROM alpine:latest

RUN apk add ca-certificates

WORKDIR /curr

COPY --from=builder /curr-build/curr .

CMD [ "./curr" ]