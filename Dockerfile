FROM golang:1.11.10
MAINTAINER wangxk xinkuanwang@gmail.com

WORKDIR /GoProject/TaibaiSupport
ADD . /GoProject/TaibaiSupport

RUN go mod tidy
RUN go build .

EXPOSE 8888

ENTRYPOINT ["./TaibaiSupport"]



