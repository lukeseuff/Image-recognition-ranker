FROM golang
EXPOSE 8080
WORKDIR /go/src/app
COPY . /go/src/app
CMD ["go", "run", "hello-world.go"]
