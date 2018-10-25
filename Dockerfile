FROM golang
EXPOSE 8080
WORKDIR /go
COPY . /go
CMD ["go", "run", "main.go"]
