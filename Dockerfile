FROM docker.io/golang:1.23.0 AS builder

WORKDIR /app

RUN apt update && apt -y dist-upgrade && apt -y install curl wget jq

ADD . .

RUN go get && go mod tidy && go build -o app.bin && chmod +x app.bin

# FROM ubuntu:latest

# COPY --from=builder /app/app.bin /app.bin 

# RUN apt update && apt -y dist-upgrade && apt -y install curl wget jq

ENTRYPOINT [ "/app/app.bin" ]
