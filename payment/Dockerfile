FROM golang:alpine AS builder

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY go.* ./

RUN go mod download

COPY . .