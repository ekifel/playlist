FROM golang:latest

RUN go version
ENV GOPATH=/
COPY ./ ./

RUN go mod download
RUN go build -o playlist ./cmd/playlist/main.go

CMD ["./playlist"]