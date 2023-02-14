FROM golang

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN GOOS=linux go build -o ./app ./cmd/app/main.go

CMD ["./app"]
