FROM golang:1.11.10 as build-env
MAINTAINER wangxk xinkuanwang@gmail.com

WORKDIR /GoProject/TaibaiSupport
ADD . /GoProject/TaibaiSupport

RUN go mod tidy
RUN go build --ldflags "-extldflags -static"


FROM alpine
RUN apk --no-cache add tzdata
WORKDIR /app
COPY --from=build-env /GoProject/TaibaiSupport/TaibaiSupport /app


EXPOSE 8888

ENTRYPOINT ["./TaibaiSupport"]



