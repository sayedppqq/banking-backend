# Build stage
FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY config.env .
COPY start.sh .
RUN chmod +x /app/start.sh

EXPOSE 8081

CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]