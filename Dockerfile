# syntax=docker/dockerfile:1
FROM golang:1.17

COPY . /mongogo
WORKDIR /mongogo
RUN go mod download

WORKDIR /mongogo/app
RUN go build -o mongogo

CMD [ "./mongogo" ]