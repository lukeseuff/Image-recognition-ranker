FROM golang
EXPOSE 80
WORKDIR /go/src/app
COPY . /go/src/app
CMD ["go", "run", "hello-world.go"]
