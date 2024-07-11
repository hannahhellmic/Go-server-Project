package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "awesomeProject/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBankAccountServiceServer
	accounts map[string]*pb.CreateResponse
	mu       sync.Mutex
}

func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.accounts[req.Name]; exists {
		return nil, fmt.Errorf("account already exists")
	}

	account := &pb.CreateResponse{
		Name:    req.Name,
		Balance: req.Balance,
	}
	s.accounts[req.Name] = account
	return account, nil
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	account, exists := s.accounts[req.Name]
	if !exists {
		return nil, fmt.Errorf("account not found")
	}
	return &pb.GetResponse{
		Name:    account.Name,
		Balance: account.Balance,
	}, nil
}

func (s *server) UpdateBalance(ctx context.Context, req *pb.UpdateBalanceRequest) (*pb.UpdateBalanceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	account, exists := s.accounts[req.Name]
	if !exists {
		return nil, fmt.Errorf("account not found")
	}

	account.Balance += req.Balance
	return &pb.UpdateBalanceResponse{
		Name:    account.Name,
		Balance: account.Balance,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.accounts[req.Name]; exists {
		delete(s.accounts, req.Name)
		return &pb.DeleteResponse{Name: req.Name}, nil
	}
	return nil, fmt.Errorf("account not found")
}

func (s *server) UpdateName(ctx context.Context, req *pb.UpdateNameRequest) (*pb.UpdateNameResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	account, exists := s.accounts[req.Name]
	if !exists {
		return nil, fmt.Errorf("account not found")
	}

	if _, exists := s.accounts[req.NewName]; exists {
		return nil, fmt.Errorf("this name is already taken")
	}

	account.Name = req.NewName
	s.accounts[req.NewName] = account
	delete(s.accounts, req.Name)

	return &pb.UpdateNameResponse{Name: req.Name, NewName: req.NewName}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5678")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBankAccountServiceServer(s, &server{
		accounts: make(map[string]*pb.CreateResponse),
	})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
