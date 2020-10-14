#base image
FROM golang:alpine as builder
#setting up working directory
WORKDIR /go/src/portexercise
#copy the source files
COPY . .
#run the build
RUN go build -o server ./cmd/server/main.go

# spins a new container
FROM alpine:latest
WORKDIR /server
# copy the binary from previous builder stage
COPY --from=builder /go/src/portexercise/server /server/
RUN chmod +x server
#run the server
CMD ["./server"]
#expose the port
EXPOSE 9000