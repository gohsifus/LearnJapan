FROM golang:1.16-alpine
RUN mkdir /LearnJapan
ADD . /LearnJapan
WORKDIR /LearnJapan
RUN mkdir /logs
RUN touch ./logs/log.txt
RUN go build -o main ./cmd
CMD ["/LearnJapan/main"]
