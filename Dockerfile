FROM golang:alpine
RUN mkdir /app 
ADD . /app/
WORKDIR /app 
RUN apk --no-cache add git
RUN go get \
    github.com/lib/pq \
    github.com/gorilla/mux \
    github.com/nlopes/slack 
RUN go build -o main .
CMD ["./main"]