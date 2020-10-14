package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"portexercise/proto"

	"google.golang.org/grpc"
)

type ProtoClient struct {
	cc proto.PortDomainClient
}

// NewProtoClient ...
func NewProtoClient(cc *grpc.ClientConn) *ProtoClient {
	service := proto.NewPortDomainClient(cc)
	return &ProtoClient{service}
}

func (pc ProtoClient) Upsert(ctx context.Context, r io.ReadCloser) error {
	dec := json.NewDecoder(r)
	// read open bracket
	_, err := dec.Token()
	if err != nil {
		log.Println("dec.Token failed to read opener", err)
		return err
	}
	var stream proto.PortDomain_UpsertClient
	for dec.More() {
		port := &proto.PortInfo{}
		key, err := dec.Token()
		if err != nil {
			log.Println("Failed to decode input stream", err)
			return err
		}
		err = dec.Decode(&port)
		if err != nil {
			log.Println("Failed to decode JSON body", err)
			return err
		}
		port.Key = key.(string)

		stream, err = pc.cc.Upsert(ctx)
		if err != nil {
			log.Println("Cannot cannot stream: ", err)
			return err
		}

		err = stream.Send(port)
		if err != nil {
			log.Println("Cannot send port info to server: ", err, stream.RecvMsg(nil))
			return err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("Cannot receive response: ", err)
		return err
	}

	log.Printf("Retured status: %s\n", res.GetStatus())
	return nil
}

func (pc ProtoClient) Fetch(ctx context.Context, key string) (io.ReadCloser, error) {
	req := &proto.ReadRequest{Identifier: key}
	stream, err := pc.cc.Read(ctx, req)
	if err != nil {
		log.Println("cannot find: ", err)
		return nil, err
	}

	var res *proto.PortInfo
	var b []byte
	for {
		res, err = stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("cannot receive response: ", err)
			return nil, err
		}

		b, err = json.Marshal(res)
		if err != nil {
			log.Printf("Json marshal failed: %s\n", err.Error())
			return nil, err
		}
		//fmt.Println(string(b))
	}

	return ioutil.NopCloser(bytes.NewReader(b)), nil
}
