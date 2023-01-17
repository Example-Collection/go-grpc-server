package main

import (
	"context"
	"errors"
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
	savedPersons []*pb.SavePersonRequest
}

func (s *personServiceServer) GetPersonInformation(_ context.Context, req *pb.GetPersonRequest) (*pb.GetPersonResponse, error) {
	log.Printf("PersonRequest(email: %v) arrived.\n", req.Email)
	// find person by email in savedPersons
	var person *pb.SavePersonRequest
	for _, p := range s.savedPersons {
		if p.Email == req.Email {
			person = p
			break
		}
	}

	if person == nil {
		return nil, errors.New("person not found")
	}

	return &pb.GetPersonResponse{
		Name:    person.Name,
		Age:     person.Age,
		Email:   person.Email,
		Message: "Successfully found!(email:" + person.Email + ")",
	}, nil
}

func (s *personServiceServer) SavePersonInformation(_ context.Context, req *pb.SavePersonRequest) (*pb.BasicMessageResponse, error) {
	log.Printf("PersonRequest(name: %v, age: %d, email: %v, password: %v) arrived.\n", req.Name, req.Age, req.Email, req.Password)
	s.savedPersons = append(s.savedPersons, req)
	return &pb.BasicMessageResponse{Message: "Successfully saved!(name:" + req.Name + ", email: " + req.Email + ")"}, nil

}

func newServer() *personServiceServer {
	savedPersons := []*pb.SavePersonRequest{
		{
			Email:    "robbyra@gmail.com",
			Age:      25,
			Name:     "sangwooAged25",
			Password: "sangwooPassword",
		},
		{
			Email:    "robbyra@gmail.com",
			Age:      26,
			Name:     "sangwooAged26",
			Password: "sangwooPassword",
		},
		{
			Email:    "robbyra@gmail.com",
			Age:      27,
			Name:     "sangwooAged27",
			Password: "sangwooPassword",
		},
		{
			Email:    "notSangwoo@gmail.com",
			Age:      1,
			Name:     "notSangwoo",
			Password: "notSangwooPassword",
		},
	}
	return &personServiceServer{savedPersons: savedPersons}
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
