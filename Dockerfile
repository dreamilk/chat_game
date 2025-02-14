FROM golang:1.23.1

WORKDIR /app

COPY . .

RUN go build -o main .

CMD ["./main server"]
