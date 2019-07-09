FROM golang:alpine

RUN apk --no-cache add git

RUN go get \
    github.com/lib/pq \
    github.com/gorilla/mux \
    github.com/nlopes/slack 

WORKDIR /app

ENTRYPOINT ["./entrypoint.sh"]