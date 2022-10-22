FROM golang:latest

ADD . .

RUN go build -o hello.exe ./hello.go

EXPOSE 80

CMD ["./hello.exe"]
