# build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . . 
RUN go build -o web_server ./cmd

# production stage
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/server .

EXPOSE 8000
CMD ["./server"]
