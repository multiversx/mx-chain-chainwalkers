FROM golang:1.13-buster

WORKDIR /chainwalkers_app

COPY . .

RUN cd ./parsing && go mod vendor

RUN cd ./parsing && go get ./... && \
    go build -o ./blocks/blocks  ./blocks && \
    go build -o ./height/height  ./height
