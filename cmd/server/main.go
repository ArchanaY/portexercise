package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"portexercise/proto"
	"portexercise/service"
	repo "portexercise/service/pkg"
	"runtime/debug"

	"google.golang.org/grpc"
)

const (
	addr = ":9000"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			errStr := fmt.Errorf("%v", e)
			debugStack := string(debug.Stack())
			log.Fatal(errStr.Error())
			log.Fatal(debugStack)
			os.Exit(1)
		}
	}()

	errChan := make(chan error)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		os.Exit(1)
	}

	store := repo.NewInMemoryStore()

	s := service.NewPortServer(store)

	grpcServer := grpc.NewServer()

	proto.RegisterPortDomainServer(grpcServer, s)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	defer func() {
		grpcServer.GracefulStop()
	}()

	select {
	case err := <-errChan:
		log.Fatalf("Fatal error: %v\n", err)
	case <-quit:
	}
}
