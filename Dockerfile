FROM golang:1.23.5

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN mkdir /app/data

RUN CGO_ENABLED=1 GOOS=linux go build -o /app/xbin /app/cmd/main.go

CMD ["/app/xbin"]
