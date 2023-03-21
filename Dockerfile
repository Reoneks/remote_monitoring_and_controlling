FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o project .

FROM scratch
WORKDIR /app
COPY --from=builder /app/project /app/project
COPY --from=builder /app/migrations/* /app/migrations/*
COPY --from=builder /app/.env /app/.env

ENTRYPOINT ["./project"]
