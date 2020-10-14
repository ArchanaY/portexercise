package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"portexercise/proto"
	"portexercise/service/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Storer interface {
	Insert(ctx context.Context, pi domain.Port) error
	Fetch(ctx context.Context, key string) (domain.Port, error)
}

type PortServer struct {
	store Storer
}

func NewPortServer(s Storer) *PortServer {
	return &PortServer{s}
}

func (ps *PortServer) Upsert(stream proto.PortDomain_UpsertServer) error {
	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		req, err := stream.Recv()

		if err == io.EOF {
			log.Println("No more data")
			break
		}

		if err != nil {
			return logError(status.Errorf(codes.Unknown, "Cannot receive port info"))
		}

		port := domain.Port{
			Key:         req.GetKey(),
			Name:        req.GetName(),
			City:        req.GetCity(),
			Country:     req.GetCountry(),
			Alias:       req.GetAlias(),
			Regions:     req.GetRegions(),
			Coordinates: req.GetCoordinates(),
			Province:    req.GetProvince(),
			Timezone:    req.GetTimezone(),
			Unlocs:      req.GetUnlocs(),
			Code:        req.GetCode(),
		}

		err = ps.store.Insert(stream.Context(), port)
		if err != nil {
			err = logError(status.Errorf(codes.Unknown, "Unable to add entry for port %v", port.Key))
			fmt.Println(err)
		}
	}

	res := &proto.UpsertReply{
		Status: "OK",
	}

	err := stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}

	return nil
}

func (ps *PortServer) Read(rr *proto.ReadRequest, stream proto.PortDomain_ReadServer) error {
	p, err := ps.store.Fetch(context.Background(), rr.GetIdentifier())
	if err != nil {
		return err
	}

	pi := &proto.PortInfo{
		Key:         p.Key,
		Name:        p.Name,
		City:        p.City,
		Country:     p.Country,
		Alias:       p.Alias,
		Regions:     p.Regions,
		Coordinates: p.Coordinates,
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocs:      p.Unlocs,
		Code:        p.Code,
	}
	err = stream.Send(pi)
	if err != nil {
		return err
	}

	return nil
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}
}

func logError(err error) error {
	if err != nil {
		log.Fatal(err)
	}
	return err
}
