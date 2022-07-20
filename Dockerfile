FROM golang:1.16-alpine
RUN apk update && apk add bash
RUN mkdir /LearnJapan
ADD . /LearnJapan
WORKDIR /LearnJapan
RUN mkdir ./logs
RUN touch ./logs/log.txt
RUN GOOS=linux GOARCH=arm64 go build -o main ./cmd
CMD ["/LearnJapan/main"]
