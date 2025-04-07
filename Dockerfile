FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /to-do-api ./internal/cmd/main.go

EXPOSE 8000

CMD [ "/to-do-api" ]