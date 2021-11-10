FROM golang:1.17-bullseye

WORKDIR /go/src/todoserver
COPY . .

# RUN go get -d -v ./...
# RUN go install -v ./...

RUN go build -o main .

EXPOSE 8888

CMD ["./main"]