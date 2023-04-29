FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o project .

ENTRYPOINT ["./project"]
