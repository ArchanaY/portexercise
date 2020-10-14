#base image
FROM golang:alpine as builder
#setting up working directory
WORKDIR /go/src/portexercise
#copy the source files
COPY . .
#run the build
RUN go build -o restclient ./cmd/client/main.go

# spins a new container
FROM alpine:latest
WORKDIR /client
# copy the binary from previous builder stage
COPY --from=builder /go/src/portexercise/restclient /client/
# change permission
RUN chmod +x restclient
#run the client
CMD ["./restclient"]
#expose the port
EXPOSE 8080