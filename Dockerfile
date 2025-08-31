# syntax=docker/dockerfile:1

FROM golang:1.21.7


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./


RUN CGO_ENABLED=0 GOOS=linux go build -o /paddle-club


EXPOSE 9090

CMD ["/paddle-club"]
