FROM golang:1.22.5-alpine AS builder

# Install necessary libraries
RUN apk update && apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 go build -o server cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
COPY config.development.json .
ENV ENV=development
CMD ["./server"]
