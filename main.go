package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
	"os"

	pb "github.com/IgorEulalio/notificationservice/proto"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

const createRepositoryTimeoutSeconds = "5"
const subject = "repository.created"

type server struct {
	pb.UnimplementedRepositoryServiceServer
	nc *nats.Conn
}

type Repository struct {
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	Visibility string `json:"visibility"`
}

func (s *server) CreateRepository(ctx context.Context, req *pb.CreateRepositoryRequest) (*pb.CreateRepositoryResponse, error) {

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
	var repository Repository
	err = json.Unmarshal(jsonData, &repository)
	if err != nil {
		return nil, err
	}

	err = createRepository(repository)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(repository)
	if err != nil {
		return nil, err
	}
	err = s.nc.Publish(subject, data)
	if err != nil {
		return nil, err
	}
	log.Printf("Created repository: Name=%s, Owner=%s, Visibility=%s", req.Name, req.Owner, req.Visibility.String())
	return &pb.CreateRepositoryResponse{Message: "Repository created successfully"}, nil
}

func createRepository(repository Repository) error {
	if repository.Visibility == "PUBLIC" {
		return fmt.Errorf("user doesn't have access to create public repositories")
	}
	log.Println("Creating repository on GitHub...")

	// Load GitHub token from environment variable
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	// Create the repository on GitHub
	url := "https://api.github.com/user/repos"
	jsonRepoData, err := json.Marshal(repository)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonRepoData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create repository: %s", resp.Status)
	}

	log.Println("Repository created on GitHub successfully!")
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("failed to establish NATS connection: %v", err)
	}
	pb.RegisterRepositoryServiceServer(s, &server{nc: nc})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
