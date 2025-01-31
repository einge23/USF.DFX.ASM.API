FROM golang:1.23.5-bookworm 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]

EXPOSE 3000
