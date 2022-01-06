package main

import (
	"context"
	pb "github.com/Example-Collection/go-grpc-server/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type personServiceServer struct {
	// 아래의 UnimplementedPersonServiceServer는 인터페이스 구현을 위해
	// 항상 embed 되어 있어야 한다.
	pb.UnimplementedPersonServiceServer

	// 데이터베이스에 사용자가 저장된 것처럼 PersonRequest들을 담는 배열을 만든다.
	savedPersons []*pb.PersonRequest
}

func personRequestToPersonResponse(req *pb.PersonRequest) *pb.PersonResponse {
	return &pb.PersonResponse{
		Email:   req.Email,
		Age:     req.Age,
		Name:    req.Name,
		Message: "Successfully saved!(name:" + req.Name + ")",
	}
}

func (s *personServiceServer) GetPersonInformation(_ context.Context, req *pb.PersonRequest) (*pb.PersonResponse, error) {
	log.Printf("PersonRequest(name: %v, age: %d, email: %v, password: %v) arrived.\n", req.Name, req.Age, req.Email, req.Password)
	s.savedPersons = append(s.savedPersons, req)
	return personRequestToPersonResponse(req), nil
}

func newServer() *personServiceServer {
	return &personServiceServer{savedPersons: []*pb.PersonRequest{}}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatal("Failed to listen on port 8081")
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPersonServiceServer(grpcServer, newServer())
	_ = grpcServer.Serve(lis)
}
