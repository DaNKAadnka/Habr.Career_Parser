FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go
RUN go build -o main ./cmd/main.go
CMD ["./main"]