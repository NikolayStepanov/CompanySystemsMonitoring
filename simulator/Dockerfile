FROM golang

WORKDIR /simulator

COPY go.mod .

RUN go mod download

COPY . .

RUN GOOS=linux go build -o ./simulator ./main.go

CMD ["./simulator"]