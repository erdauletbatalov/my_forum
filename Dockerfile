FROM golang:1.17.6
RUN mkdir /app
WORKDIR /app
ADD . /app
RUN go build -o server ./cmd/
ENV TZ="Asia/Almaty"
CMD ["/app/server"]
