FROM golang:alpine
RUN mkdir /app
COPY src /app/
WORKDIR /app
RUN go build -o main .
RUN adduser -S -D -H -h /app appuser
USER appuser
ENV VERSION v1.0
CMD ["./main"]
