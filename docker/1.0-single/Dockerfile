FROM golang:alpine as builder
RUN mkdir /app
ADD src /app/
WORKDIR /app
RUN go build -o main .
RUN adduser -S -D -H -h /app appuser
USER appuser
CMD ["./main"]
