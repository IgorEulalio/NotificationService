package server

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IgorEulalio/notificationservice/biz"
	"github.com/IgorEulalio/notificationservice/model"
	pb "github.com/IgorEulalio/notificationservice/proto"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/encoding/protojson"
)

const subject = "repository.created"

type Server struct {
	pb.UnimplementedRepositoryServiceServer
	Nc *nats.Conn
}

func (s *Server) CreateRepository(ctx context.Context, req *pb.CreateRepositoryRequest) (*pb.CreateRepositoryResponse, error) {

	repoProto := &pb.Repository{
		Name:       req.Name,
		Owner:      req.Owner,
		Visibility: req.Visibility,
	}
	jsonData, err := protojson.Marshal(repoProto)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into the Go struct
	var repository model.Repository
	err = json.Unmarshal(jsonData, &repository)
	if err != nil {
		return nil, err
	}

	err = biz.CreateRepository(repository)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(repository)
	if err != nil {
		return nil, err
	}
	err = s.Nc.Publish(subject, data)
	if err != nil {
		return nil, err
	}
	log.Printf("Created repository: Name=%s, Owner=%s, Visibility=%s", req.Name, req.Owner, req.Visibility.String())
	return &pb.CreateRepositoryResponse{Message: "Repository created successfully"}, nil
}
