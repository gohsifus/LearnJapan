FROM golang:1.16-alpine
RUN mkdir /LearnJapan
ADD . /LearnJapan
WORKDIR /LearnJapan
RUN go build -o main ./cmd
CMD ["/LearnJapan/main"]