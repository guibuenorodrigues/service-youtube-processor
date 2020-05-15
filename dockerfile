FROM golang:latest

LABEL version="1.0"

ENV GO111MODULE=on

WORKDIR /app/build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/build .

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# Expose port to connect to rabbit
#EXPOSE 5672

CMD ["./out/build"]