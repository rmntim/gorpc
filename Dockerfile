FROM golang:alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o server /app/cmd/sso/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/config/local.yaml /app/config.yaml
COPY --from=builder /app/migrations/ /app/migrations/

EXPOSE 6969

CMD ["/app/server"]
