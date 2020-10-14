package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"os/signal"
	"portexercise/client"
	"runtime/debug"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

const (
	serverAddr = "server:9000"
	localAddr  = ":9000"
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
		os.Exit(1)
	}
	defer conn.Close()

	pc := client.NewProtoClient(conn)

	rf := client.NewFramework(pc)

	router := mux.NewRouter()

	router.Methods("GET").Path("/heartbeat").HandlerFunc(rf.RespondToHeartBeat)
	router.Methods("POST").Path("/port").HandlerFunc(rf.PostHandler)
	router.Methods("GET").Path("/port/{port}").HandlerFunc(rf.GetHanlder)

	go func() {
		log.Println("client starting")
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.Fatalf("listenAndServe failed: %v", err)
		}
	}()

	<-quit

	log.Println("client stopped")
}
