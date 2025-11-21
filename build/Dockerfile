FROM golang:1.25-alpine as builder

WORKDIR /app
ADD go.mod .
ADD go.sum .
RUN go mod download

COPY cmd cmd
COPY internal internal

RUN GOOS=linux go build -o main ./cmd


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
ADD configs/config.json configs/config.json

EXPOSE 8080

CMD ["./main"]