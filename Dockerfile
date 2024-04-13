FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./main.go

RUN chmod +x /app/start.sh


# Run the binary when the container starts

CMD ["./app"]



