FROM golang:alpine as builder

WORKDIR /curr-build

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init -g cmd/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o curr cmd/main.go

FROM alpine:latest

RUN apk add ca-certificates

WORKDIR /curr

COPY --from=builder /curr-build/curr .

EXPOSE 2000

CMD [ "./curr" ]